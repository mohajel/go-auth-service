package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

type mongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(db *mongo.Database) Repository {
	return &mongoRepo{
		collection: db.Collection("users"),
	}
}

func (r *mongoRepo) Create(ctx context.Context, u *User) error {
	_, err := r.collection.InsertOne(ctx, u)
	return err
}

func (r *mongoRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoRepo) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
