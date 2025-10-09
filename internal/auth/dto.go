package auth

type ProfileResponse struct {
	ID       string
	Username string
	Email    string
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}
