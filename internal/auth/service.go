package auth

import (
	"context"
	"errors"

	"go-auth-service/pkg/logger"

	"go-auth-service/internal/user"

	userRedis "go-auth-service/internal/redis"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
}

type SimpleTokenManager struct{}

func (t *SimpleTokenManager) GenerateAccessToken(userID string) (string, error) {
	// TODO: Replace with real JWT generation
	return "access-token-" + userID, nil
}

func (t *SimpleTokenManager) GenerateRefreshToken(userID string) (string, error) {
	// TODO: Replace with real JWT generation
	return "refresh-token-" + userID + "-" + uuid.NewString(), nil
}

type NatsConn interface {
	Publish(subject string, data []byte) error
}
type Service struct {
	userRepo  user.Repository
	tokenMgr  TokenManager
	natsConn  NatsConn
	redisRepo *userRedis.Repository
}

func NewService(
	userRepo user.Repository,
	tokenMgr TokenManager,
	natsConn NatsConn,
	redisRepo *userRedis.Repository,
) *Service {
	return &Service{
		userRepo:  userRepo,
		tokenMgr:  tokenMgr,
		natsConn:  natsConn,
		redisRepo: redisRepo,
	}
}

// Helper Functions
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // TODO: Check this hash method
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateID() string {
	return uuid.NewString()
}

// Service
func (s *Service) RegisterUser(req *RegisterRequest) error {
	ctx := context.Background()

	existingUser, _ := s.userRepo.FindByUsername(ctx, req.Username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	newUser := &user.User{
		ID:           generateID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return err
	}

	// Publish event to NATS
	if s.natsConn != nil {
		_ = s.natsConn.Publish("user.registered", []byte(newUser.ID))
	}

	logger.Info("User registered: " + newUser.ID)
	return nil
}

func (s *Service) LoginUser(req *LoginRequest) (*LoginResponse, error) {
	ctx := context.Background()

	u, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil || u == nil {
		return nil, errors.New("invalid credentials")
	}

	if !checkPasswordHash(req.Password, u.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.tokenMgr.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenMgr.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}
	// Store refresh roken in Redis
	if err := s.redisRepo.StoreRefreshToken(u.ID, refreshToken); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Validates refresh token and generates new one
func (s *Service) Refresh(refreshToken string) (*LoginResponse, error) {
	userID, err := s.redisRepo.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, err := s.tokenMgr.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenMgr.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	// Replace old refresh token
	if err := s.redisRepo.DeleteRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	if err := s.redisRepo.StoreRefreshToken(userID, newRefreshToken); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(refreshToken string) error {
	userID, err := s.redisRepo.ValidateRefreshToken(refreshToken)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	if err := s.redisRepo.DeleteRefreshToken(refreshToken); err != nil {
		return err
	}

	if s.natsConn != nil {
		_ = s.natsConn.Publish("user.logged_out", []byte(userID))
	}

	logger.Info("User logged out: " + userID)
	return nil
}

func (s *Service) Profile(userID string) (*ProfileResponse, error) {
	u, err := s.userRepo.FindByID(context.Background(), userID)
	if err != nil || u == nil {
		return nil, errors.New("user not found")
	}

	return &ProfileResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}
