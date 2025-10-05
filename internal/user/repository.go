package user

import "context"

type Repository interface {
    Create(ctx context.Context, u *User) error
    FindByUsername(ctx context.Context, username string) (*User, error)
    FindByID(ctx context.Context, id string) (*User, error)
}

func NewMongoRepo(db interface{}) Repository {
    // TODO: Return a new MongoDB repository instance
    return nil
}
