package user

type User struct {
	ID       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"` // TODO: Add email field
	Password string `bson:"password" json:"password"`
	// TODO: Optional RefreshToken
}
