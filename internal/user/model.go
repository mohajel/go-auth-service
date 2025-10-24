package user

import (
	"go-auth-service/pkg/hash"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

func NewUser(username, password, email string) (*User, error) {
	hashedPassword, err := hash.Generate(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       primitive.NewObjectID().Hex(),
		Username: username,
		Password: hashedPassword,
		Email:    email,
	}, nil
}

func (u *User) VerifyPassword(password string) bool {
	return hash.Verify(password, u.Password)
}
