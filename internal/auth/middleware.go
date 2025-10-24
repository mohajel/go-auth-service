package auth

import (
	"strings"

	"go-auth-service/internal/token"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "user_id"
)

func JWTMiddleware(tokenManager *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := tokenManager.Validate(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
