package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

type AuthMiddleware struct {
	Decrypter domain.Decrypter
}

func NewAuthMiddleware(decrypter domain.Decrypter) *AuthMiddleware {
	return &AuthMiddleware{
		Decrypter: decrypter,
	}
}

func (h *AuthMiddleware) Handle(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if len(strings.Split(token, "Bearer")) < 2 {
		c.JSON(401, outputs.ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	data, err := h.Decrypter.Decrypt(strings.TrimSpace(strings.Split(token, "Bearer")[1]))

	if err != nil {
		c.JSON(401, outputs.ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	c.Set("currentCustomer", data)
}
