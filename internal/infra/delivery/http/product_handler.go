package http

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	"github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

type productHandler struct {
	getProductUc  domain.GetProductUseCase
	listProductUc domain.ListProductsUseCase
}

func NewProductHandler(
	r *gin.RouterGroup,
	getProductUc domain.GetProductUseCase,
	listProductUc domain.ListProductsUseCase,
) *gin.RouterGroup {
	handler := &productHandler{
		getProductUc:  getProductUc,
		listProductUc: listProductUc,
	}

	productRoutes := r.Group("/products")
	productRoutes.GET("/:productId", handler.GetProduct)
	productRoutes.GET("/", handler.ListProducts)

	return productRoutes
}

// GetProduct godoc
// @Summary Get product details by ID
// @Description Retrieves detailed information about a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} outputs.ErrorResponse "Invalid input"
// @Failure 404 {object} outputs.ErrorResponse "Product not found"
// @Failure 500 {object} outputs.ErrorResponse "Internal server error"
// @Router /api/products/{productId} [get]
func (h *productHandler) GetProduct(c *gin.Context) {
	pid := c.Param("productId")

	if pid == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		c.Next()
		return
	}

	ps, err := h.getProductUc.Execute(c.Request.Context(), pid)

	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, ps)
	return
}

// ListProducts godoc
// @Summary List all products with pagination
// @Description Returns a paginated list of products
// @Tags products
// @Accept json
// @Produce json
// @Param page query string false "Page number"
// @Param size query string false "Page size (default: 20)"
// @Success 200 {array} domain.Product
// @Failure 400 {object} outputs.ErrorResponse "Invalid input"
// @Failure 500 {object} outputs.ErrorResponse "Internal server error"
// @Router /api/products [get]
func (h *productHandler) ListProducts(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("size")

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 0
	}

	sizeNum, err := strconv.Atoi(size)
	if err != nil {
		sizeNum = 20
	}

	if pageNum < 0 || sizeNum < 0 {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		c.Next()
		return
	}

	products, err := h.listProductUc.Execute(c.Request.Context(), sizeNum, pageNum*sizeNum)

	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(200, []struct{}{}) // not found on a list route should be an empty array
			return
		}
		fmt.Printf("error fetching products: %v\n", err)
		c.JSON(500, outputs.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.JSON(200, products)
	return
}
