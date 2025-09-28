package token

func NewManager(secret string) *Manager {
    // TODO: Return a new JWT manager instance
    return nil
}

type Manager struct {
    // TODO: secret and methods for JWT generation/validation
}

func (m *Manager) Generate(userID string) string {
    // TODO: Generate access token
    return ""
}

func (m *Manager) Validate(token string) (string, error) {
    // TODO: Validate token and return userID
    return "", nil
}
