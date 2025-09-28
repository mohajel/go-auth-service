package auth

import "github.com/gin-gonic/gin"

type Handler struct {
    // TODO: Auth service
}

func NewHandler(s *Service) *Handler {
    // TODO: Return a new Handler instance
    return nil
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
    // TODO: Mount routes: /register, /login, /refresh, /logout, /me/profile
}

func (h *Handler) Register(c *gin.Context) {}
func (h *Handler) Login(c *gin.Context) {}
func (h *Handler) Refresh(c *gin.Context) {}
func (h *Handler) Logout(c *gin.Context) {}
func (h *Handler) Profile(c *gin.Context) {}
