package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTokenManager_CreateAndVerifyToken(t *testing.T) {
	mgr := NewTokenManager("test-secret")

	token, err := mgr.CreateToken("123", "alice", time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := mgr.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "123", payload.UserID)
	assert.Equal(t, "alice", payload.Username)
}

func TestTokenManager_InvalidToken(t *testing.T) {
	mgr := NewTokenManager("test-secret")

	_, err := mgr.VerifyToken("wrong.token.here")
	assert.Error(t, err)
}

func TestJWTMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := NewTokenManager("test-secret")

	token, _ := mgr.CreateToken("123", "alice", time.Minute)

	r := gin.New()
	r.Use(JWTMiddleware(mgr))
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "username": username})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "alice")
}

func TestJWTMiddleware_InvalidHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := NewTokenManager("test-secret")

	r := gin.New()
	r.Use(JWTMiddleware(mgr))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
