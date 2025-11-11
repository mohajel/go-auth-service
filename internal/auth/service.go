package auth

import (
	"context"
	"encoding/json"
	"fmt"
	appnats "go-auth-service/internal/nats"
	"go-auth-service/internal/redis"
	"go-auth-service/internal/token"
	"go-auth-service/internal/user"

	"os"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service struct {
	userRepo     user.Repository
	tokenManager *token.Manager
	natsConn     *nats.Conn
	redisRepo    *redis.Repository
	oauthConfig  *oauth2.Config
}

func NewService(
	userRepo user.Repository,
	tokenManager *token.Manager,
	natsConn *nats.Conn,
	redisRepo *redis.Repository,
) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		natsConn:     natsConn,
		redisRepo:    redisRepo,
		oauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	newUser, err := user.NewUser(req.Username, req.Password, req.Email)
	if err != nil {
		return err
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return err
	}

	js, err := s.natsConn.JetStream()
	if err != nil {
		return err
	}

	appnats.PublishUserRegistered(js, newUser.ID)
	return nil
}

func (s *Service) Login(username, password string) (string, string, error) {
	u, err := s.userRepo.FindByUsername(context.Background(), username)
	if err != nil {
		return "", "", err
	}

	if !u.VerifyPassword(password) {
		return "", "", fmt.Errorf("invalid credentials")
	}

	accessToken := s.tokenManager.Generate(u.ID)
	refreshToken := token.GenerateRefreshToken()

	if err := s.redisRepo.StoreRefreshToken(context.Background(), refreshToken, u.ID); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(refreshToken string) (string, string, error) {
	userID, err := s.redisRepo.GetUserIDByRefreshToken(context.Background(), refreshToken)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	newAccessToken := s.tokenManager.Generate(userID)
	newRefreshToken := token.GenerateRefreshToken()

	ctx := context.Background()
	if err := s.redisRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return "", "", err
	}
	if err := s.redisRepo.StoreRefreshToken(ctx, newRefreshToken, userID); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *Service) Logout(refreshToken string) error {
	return s.redisRepo.DeleteRefreshToken(context.Background(), refreshToken)
}

func (s *Service) Profile(userID string) (*ProfileResponse, error) {
	u, err := s.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return &ProfileResponse{
		ID:       u.ID,
		Username: u.Username,
	}, nil
}

func (s *Service) GoogleLogin(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *Service) GoogleCallback(code string) (string, string, error) {
	t, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", "", err
	}

	client := s.oauthConfig.Client(context.Background(), t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	userInfo := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", "", err
	}

	email := userInfo["email"].(string)

	// بررسی اینکه کاربر قبلا وجود داشته یا نه
	u, err := s.userRepo.FindByEmail(context.Background(), email)
	if err != nil {
		// اگر نبود، کاربر جدید بساز
		newUser, _ := user.NewUser(uuid.New().String(), "", email)
		if err := s.userRepo.Create(context.Background(), newUser); err != nil {
			return "", "", err
		}
		u = newUser
	}

	// ساخت JWT و Refresh Token
	accessToken := s.tokenManager.Generate(u.ID)
	refreshToken := token.GenerateRefreshToken()
	if err := s.redisRepo.StoreRefreshToken(context.Background(), refreshToken, u.ID); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
