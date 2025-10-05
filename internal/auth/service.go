package auth

type Service struct {
    // TODO: userRepo, tokenManager, nats, redisRepo
}

type dummy interface{}

func NewService(
    userRepo dummy,
    tokenManager dummy,
    natsConn dummy,
    redisRepo dummy,
) *Service {
    // TODO: Return a new Auth service instance
    return &Service{}
}

func (s *Service) Register(username, password string) error {
    // TODO: Store user and publish event
    return nil
}

func (s *Service) Login(username, password string) (string, string, error) {
    // TODO: Verify password and create access + refresh token
    return "", "", nil
}

func (s *Service) Refresh(refreshToken string) (string, string, error) {
    // TODO: Verify refresh token and generate new token
    return "", "", nil
}

func (s *Service) Logout(refreshToken string) error {
    // TODO: Invalidate refresh token and publish event
    return nil
}

func (s *Service) Profile(userID string) (*ProfileResponse, error) {
    // TODO: Return user profile data
    return nil, nil
}
