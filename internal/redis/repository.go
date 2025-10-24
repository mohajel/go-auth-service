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

func (r *Repository) StoreRefreshToken(ctx context.Context, token, userID string) error {
    return r.client.Set(ctx, "refresh_token:"+token, userID, 7*24*time.Hour).Err()
}

func (r *Repository) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
    return r.client.Get(ctx, "refresh_token:"+token).Result()
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, token string) error {
    return r.client.Del(ctx, "refresh_token:"+token).Err()
}