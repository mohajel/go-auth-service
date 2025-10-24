package auth

import (
	"context"
	"fmt"
	appnats "go-auth-service/internal/nats"
	"go-auth-service/internal/redis"
	"go-auth-service/internal/token"
	"go-auth-service/internal/user"

	"github.com/nats-io/nats.go"
)

type Service struct {
	userRepo     user.Repository
	tokenManager *token.Manager
	natsConn     *nats.Conn
	redisRepo    *redis.Repository
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
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	// Create new user
	newUser, err := user.NewUser(req.Username, req.Password, req.Email)
	if err != nil {
		return err
	}

	// Store in database
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return err
	}

	// Publish event
	js, err := s.natsConn.JetStream()
	if err != nil {
		return err
	}

	appnats.PublishUserRegistered(js, newUser.ID)
	return nil
}

func (s *Service) Login(username, password string) (string, string, error) {
	// Find user
	u, err := s.userRepo.FindByUsername(context.Background(), username)
	if err != nil {
		return "", "", err
	}

	// Verify password
	if !u.VerifyPassword(password) {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken := s.tokenManager.Generate(u.ID)
	refreshToken := token.GenerateRefreshToken()

	// Store refresh token
	if err := s.redisRepo.StoreRefreshToken(context.Background(), refreshToken, u.ID); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(refreshToken string) (string, string, error) {
	// Get user ID from refresh token
	userID, err := s.redisRepo.GetUserIDByRefreshToken(context.Background(), refreshToken)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	newAccessToken := s.tokenManager.Generate(userID)
	newRefreshToken := token.GenerateRefreshToken()

	// Store new refresh token and delete old one
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
