package http

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/presentation/inputs"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

type WishlistHandler struct {
	createWishlistUseCase domain.CreateWishlistUseCase
	renameWishlistUseCase domain.UpdateWishlistNameUseCase
	deleteWishlistUseCase domain.DeleteWishlistUseCase
	getWishlistUseCase    domain.ShowWishlistUseCase
	productAdder          domain.AddProductToWishlistUseCase
}

func SetupWishlistHandler(
	r *gin.RouterGroup,
	auth gin.HandlerFunc,
	createWishlistUseCase domain.CreateWishlistUseCase,
	renameWishlistUseCase domain.UpdateWishlistNameUseCase,
	deleteWishlistUseCase domain.DeleteWishlistUseCase,
	getWishlistUseCase domain.ShowWishlistUseCase,
	productAdder domain.AddProductToWishlistUseCase,

) {
	handler := &WishlistHandler{
		createWishlistUseCase: createWishlistUseCase,
		renameWishlistUseCase: renameWishlistUseCase,
		deleteWishlistUseCase: deleteWishlistUseCase,
		getWishlistUseCase:    getWishlistUseCase,
		productAdder:          productAdder,
	}

	wishlistRoutes := r.Group("/:customerId/wishlists")
	wishlistRoutes.Use(auth)
	wishlistRoutes.POST("/", handler.CreateWishList)
	wishlistRoutes.PATCH("/:wishListId", handler.RenameWishlist)
	wishlistRoutes.DELETE("/:wishListId", handler.DeleteWishlist)
	wishlistRoutes.GET("/:wishListId", handler.GetWishlist)
	wishlistRoutes.POST("/:wishListId/products", handler.AddProductToWishlist)

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

	wishlistId, err := h.createWishlistUseCase.CreateWishlist(c.Request.Context(), currentCustomer.ID, cid, input.Title)

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

// RenameWishlist godoc
// @Summary Renames an existing wishlist
// @Tags wishlists
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Param wishListId path string true "Wishlist ID"
// @Param wishlist body inputs.CreateWishlistInput true "Wishlist data"
// @Success 200 {object} outputs.CreateWishlistResponse
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 404 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId}/wishlists/{wishListId} [patch]
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func (h WishlistHandler) RenameWishlist(c *gin.Context) {
	cid := c.Param("customerId")
	wid := c.Param("wishListId")

	if cid == "" || wid == "" {
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

	wishlistId, err := h.renameWishlistUseCase.RenameWishlist(c.Request.Context(), currentCustomer.ID, cid, wid, input.Title)

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

		fmt.Println(err)
		c.JSON(500, outputs.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}
	c.JSON(200, outputs.CreateWishlistResponse{
		ID: wishlistId,
	})
	return
}

// DeleteWishlist godoc
// @Summary Deletes an existing wishlist
// @Tags wishlists
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Param wishListId path string true "Wishlist ID"
// @Success 204
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 404 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId}/wishlists/{wishListId} [delete]
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func (h WishlistHandler) DeleteWishlist(c *gin.Context) {
	cid := c.Param("customerId")
	wid := c.Param("wishListId")

	if cid == "" || wid == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	var currentCustomer domain.OutgoingCustomer
	str, _ := c.Get("currentCustomer")
	json.Unmarshal([]byte(str.(string)), &currentCustomer)

	err := h.deleteWishlistUseCase.DeleteWishlist(c.Request.Context(), currentCustomer.ID, cid, wid)

	if err != nil {
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

		fmt.Println(err)
		c.JSON(500, outputs.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}
	c.JSON(204, gin.H{})
	return
}

// GetWishlist godoc
// @Summary Retrieves an existing wishlist
// @Tags wishlists
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Param wishListId path string true "Wishlist ID"
// @Success 200 {object} domain.FullfilledWishlist
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 404 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId}/wishlists/{wishListId} [get]
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func (h WishlistHandler) GetWishlist(c *gin.Context) {
	cid := c.Param("customerId")
	wid := c.Param("wishListId")

	if cid == "" || wid == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	var currentCustomer domain.OutgoingCustomer
	str, _ := c.Get("currentCustomer")
	json.Unmarshal([]byte(str.(string)), &currentCustomer)

	list, err := h.getWishlistUseCase.ShowWishlist(c.Request.Context(), currentCustomer.ID, cid, wid)

	if err != nil {
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

		fmt.Println(err)
		c.JSON(500, outputs.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}
	c.JSON(200, list)
	return
}

// AddProductToWishlist godoc
// @Summary Add product to wishlist
// @Tags wishlists
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Param wishListId path string true "Wishlist ID"
// @Param wishlist body inputs.ProdcutToWishlistInput true "product data"
// @Success 204
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 404 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId}/wishlists/{wishListId}/products [post]
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func (h WishlistHandler) AddProductToWishlist(c *gin.Context) {
	h.ensureParams(c)
	customer := h.getCustomerFromContext(c)

	if customer == nil {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Missing data"})
		return
	}

	var data inputs.ProdcutToWishlistInput
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := h.productAdder.AddProduct(c.Request.Context(), customer.ID, c.Param("customerId"), c.Param("wishListId"), data.ProductId)

	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(204, gin.H{})
	return
}

func (h WishlistHandler) ensureParams(c *gin.Context) {
	cid := c.Param("customerId")
	wid := c.Param("wishListId")

	if cid == "" || wid == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}
}

func (h WishlistHandler) getCustomerFromContext(c *gin.Context) *domain.OutgoingCustomer {
	var currentCustomer domain.OutgoingCustomer
	str, e := c.Get("currentCustomer")

	if !e {
		return nil
	}

	err := json.Unmarshal([]byte(str.(string)), &currentCustomer)
	if err != nil {
		return nil
	}

	return &currentCustomer
}
