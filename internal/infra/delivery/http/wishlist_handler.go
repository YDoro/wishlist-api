package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	"github.com/ydoro/wishlist/internal/presentation/inputs"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

type WishlistHandler struct {
	createWishlistUseCase domain.CreateWishlistUseCase
	deleteWishlistUseCase domain.DeleteWishlistUseCase
	getWishlistUseCase    domain.ShowWishlistUseCase
	updateWishlistUsecase domain.UpdateWishListUseCase
}

func SetupWishlistHandler(
	r *gin.RouterGroup,
	auth gin.HandlerFunc,
	createWishlistUseCase domain.CreateWishlistUseCase,
	deleteWishlistUseCase domain.DeleteWishlistUseCase,
	getWishlistUseCase domain.ShowWishlistUseCase,
	updateWishlistUsecase domain.UpdateWishListUseCase,
) {
	handler := &WishlistHandler{
		createWishlistUseCase: createWishlistUseCase,
		deleteWishlistUseCase: deleteWishlistUseCase,
		getWishlistUseCase:    getWishlistUseCase,
		updateWishlistUsecase: updateWishlistUsecase,
	}

	wishlistRoutes := r.Group("/:customerId/wishlists")
	wishlistRoutes.Use(auth)
	wishlistRoutes.POST("/", handler.CreateWishList)
	wishlistRoutes.PUT("/:wishListId", handler.UpdateWishlist)
	wishlistRoutes.PATCH("/:wishListId", handler.UpdateWishlist)
	wishlistRoutes.DELETE("/:wishListId", handler.DeleteWishlist)
	wishlistRoutes.GET("/:wishListId", handler.GetWishlist)

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

	currentCustomer := GetCustomerFromContext(c)
	wishlistId, err := h.createWishlistUseCase.CreateWishlist(c.Request.Context(), currentCustomer.ID, cid, input.Title)

	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, outputs.CreateWishlistResponse{
		ID: wishlistId,
	})
	return
}

// UpdateWishlist godoc
// @Summary update wishlist
// @Description update wishlist, updates both title and items if given, if you want to add a single item you need to pass the whole wishlist
// @Tags wishlists
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Param wishListId path string true "Wishlist ID"
// @Param wishlist body inputs.UpdateWishlistInput true "Wishlist data"
// @Success 204
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 404 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId}/wishlists/{wishListId} [put]
// @Router /api/customers/{customerId}/wishlists/{wishListId} [patch]
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func (h WishlistHandler) UpdateWishlist(c *gin.Context) {
	h.ensureParams(c)

	var input inputs.UpdateWishlistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	currentCustomer := GetCustomerFromContext(c)

	wl := &domain.Wishlist{
		Title:      input.Title,
		ID:         c.Param("wishListId"),
		CustomerId: c.Param("customerId"),
		Items:      input.Items,
	}

	err := h.updateWishlistUsecase.UpdateWishlist(c.Request.Context(), currentCustomer.ID, wl)

	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(204, gin.H{})
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
	h.ensureParams(c)
	currentCustomer := GetCustomerFromContext(c)
	err := h.deleteWishlistUseCase.DeleteWishlist(c.Request.Context(), currentCustomer.ID, c.Param("customerId"), c.Param("wishListId"))

	if err != nil {
		HandleError(c, err)
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
	h.ensureParams(c)
	currentCustomer := GetCustomerFromContext(c)
	list, err := h.getWishlistUseCase.ShowWishlist(c.Request.Context(), currentCustomer.ID, c.Param("customerId"), c.Param("wishListId"))

	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, list)
	return
}

func (h WishlistHandler) ensureParams(c *gin.Context) {
	cid := c.Param("customerId")
	wid := c.Param("wishListId")

	if cid == "" || wid == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		c.Next()
		return
	}
}
