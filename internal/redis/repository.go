package redis

func NewRepository(r interface{}) *Repository {
    // TODO: Return new Redis repository instance
    return nil
}

type Repository struct {
    // TODO: Redis connection
}

func (r *Repository) StoreRefreshToken(userID, token string) error {
    // TODO: Store refresh token with TTL
    return nil
}

func (r *Repository) ValidateRefreshToken(token string) (string, error) {
    // TODO: Validate refresh token
    return "", nil
}

func (r *Repository) DeleteRefreshToken(token string) error {
    // TODO: Delete refresh token
    return nil
}
