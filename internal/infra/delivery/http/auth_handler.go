package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	"github.com/ydoro/wishlist/internal/presentation/inputs"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

type AuthHandler struct {
	AuthUseCase domain.Authenticator
}

// PasswordAuthentication godoc
// @Summary Authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param customer body inputs.PwdAuth true "user credentials"
// @Success 200 {object} outputs.AuthSuccessResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) PasswordAuthentication(c *gin.Context) {
	var credentials inputs.PwdAuth

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.AuthUseCase.Authenticate(c, credentials)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, outputs.AuthSuccessResponse{Token: token})
}
func NewAuthHandler(r *gin.RouterGroup, uc domain.Authenticator) {
	handler := &AuthHandler{
		AuthUseCase: uc,
	}

	authRoutes := r.Group("/auth")
	authRoutes.POST("/login", handler.PasswordAuthentication)
}
