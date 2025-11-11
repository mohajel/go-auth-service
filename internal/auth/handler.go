package auth

import (
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/refresh", h.Refresh)
	r.POST("/logout", h.Logout)
	r.GET("/me/profile", h.Profile)
	r.GET("/auth/google/login", h.GoogleLogin)
	r.GET("/auth/google/callback", h.GoogleCallback)
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newAccess, newRefresh, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, TokenResponse{AccessToken: newAccess, RefreshToken: newRefresh})
}

func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Logout(req.RefreshToken); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Logged out"})
}

func (h *Handler) Profile(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(401, gin.H{"error": "missing authorization header"})
		return
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(401, gin.H{"error": "invalid authorization header"})
		return
	}

	tokenString := parts[1]

	userID, err := h.service.tokenManager.Validate(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.Profile(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, profile)
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	state := uuid.New().String()
	url := h.service.GoogleLogin(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not provided"})
		return
	}

	accessToken, refreshToken, err := h.service.GoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	redirectURL := "http://127.0.0.1:5500/home?access=" + accessToken + "&refresh=" + refreshToken // Adjust the URL as needed
	c.Redirect(http.StatusSeeOther, redirectURL)
}
