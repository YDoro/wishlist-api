package http

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/presentation/inputs"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

type WishlistHandler struct {
	// Use cases
	CreateWishlistUseCase domain.CreateWishlistUseCase
}

func SetupWishlistHandler(
	r *gin.RouterGroup,
	auth gin.HandlerFunc,
	createWishlistUseCase domain.CreateWishlistUseCase,
) {
	handler := &WishlistHandler{
		CreateWishlistUseCase: createWishlistUseCase,
	}

	wishlistRoutes := r.Group("/:customerId/wishlists")
	wishlistRoutes.Use(auth)
	wishlistRoutes.POST("/", handler.CreateWishList)
}

// CreateWishList godoc
// @Summary Creates a new wishlist
// @Tags wishlists
// @Accept json
// @Produce json
// @Param customer path string true "Customer ID"
// @Param wishlist body inputs.CreateWishlistInput true "Wishlist data"
// @Success 201 {object} outputs.CreateWishlistResponse
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId}/wishlists [post]
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func (h WishlistHandler) CreateWishList(c *gin.Context) {
	cid := c.Param("customerId")
	if cid == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	var input inputs.CreateWishlistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var currentCustomer domain.OutgoingCustomer
	str, _ := c.Get("currentCustomer")
	json.Unmarshal([]byte(str.(string)), &currentCustomer)

	wishlistId, err := h.CreateWishlistUseCase.CreateWishlist(c.Request.Context(), currentCustomer.ID, cid, input.Title)

	if err != nil {
		if _, ok := err.(*e.ValidationError); ok {
			c.JSON(400, outputs.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		if e.IsNotFoundError(err) {
			c.JSON(404, outputs.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		if e.IsUnauthorizedError(err) {
			c.JSON(401, outputs.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(500, outputs.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}
	c.JSON(201, outputs.CreateWishlistResponse{
		ID: wishlistId,
	})
	return
}
