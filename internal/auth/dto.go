package auth

type RegisterRequest struct {
    Username string
    Password string
	Email    string
}

type LoginRequest struct {
    Username string
    Password string
	Email    string
}

type ProfileResponse struct {
    ID       string
    Username string
}

type RefreshRequest struct {
    RefreshToken string
}

type TokenResponse struct {
    AccessToken  string
    RefreshToken string
}
