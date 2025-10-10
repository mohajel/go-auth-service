// main.go
package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret string
	alg    string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{
		secret: secret,
		alg:    jwt.SigningMethodHS256.Alg(),
	}
}

func (m *TokenManager) CreateToken(userID, username string, duration time.Duration) (string, error) {
	claims := TokenPayload{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *TokenManager) VerifyToken(tokenStr string) (*TokenPayload, error) {
	claims := &TokenPayload{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != m.alg {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func JWTMiddleware(mgr *TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.Fields(authHeader) // split by spaces
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be: Bearer <token>"})
			return
		}
		tokenStr := parts[1]

		payload, err := mgr.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token", "detail": err.Error()})
			return
		}

		c.Set("user_id", payload.UserID)
		c.Set("username", payload.Username)

		c.Next()
	}
}

func ProtectedHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"message":  "you are authorized",
		"user_id":  userID,
		"username": username,
	})
}

func main() {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		secret = "very-secret-testing-key"
	}

	mgr := NewTokenManager(secret)
	sampleToken, err := mgr.CreateToken("123", "alice", time.Hour)
	if err != nil {
		panic(err)
	}
	fmt.Println("sample token (use this for curl):")
	fmt.Println(sampleToken)

	r := gin.Default()

	protected := r.Group("/protected")
	protected.Use(JWTMiddleware(mgr))
	{
		protected.GET("/", ProtectedHandler)
	}

	fmt.Println("server running on :8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
