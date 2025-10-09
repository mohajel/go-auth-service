package user

type User struct {
	ID           string `bson:"id" json:"id"`
	Username     string `bson:"username" json:"username"`
	Email        string `bson:"email" json:"email"` // TODO: Add email field
	PasswordHash string `bson:"password_hash" json:"-"`
	// TODO: Optional RefreshToken
}
