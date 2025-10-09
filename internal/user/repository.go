package user

import (
	"context"
	"errors"
	"go-auth-service/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	collection := db.Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		logger.Error("failed to create user indexes: " + err.Error())
	}

	return &mongoRepo{
		collection: collection,
	}
}

func (r *mongoRepo) Create(ctx context.Context, u *User) error {
	if u == nil {
		return errors.New("user is nil")
	}
	_, err := r.collection.InsertOne(ctx, u)
	return err
}

func (r *mongoRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *mongoRepo) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}
