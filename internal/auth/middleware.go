package auth

import "github.com/gin-gonic/gin"

func JWTMiddleware() gin.HandlerFunc {
    // TODO: Validate access token
    return func(c *gin.Context) {
        c.Next()
    }
}
