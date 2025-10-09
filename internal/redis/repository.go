package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) StoreRefreshToken(userID, token string) error {
	ctx := context.Background()
	return r.client.Set(ctx, token, userID, 7*24*time.Hour).Err()
}

func (r *Repository) ValidateRefreshToken(token string) (string, error) {
	ctx := context.Background()
	userID, err := r.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return "", nil // token not found
	}
	return userID, err
}

func (r *Repository) DeleteRefreshToken(token string) error {
	ctx := context.Background()
	return r.client.Del(ctx, token).Err()
}
