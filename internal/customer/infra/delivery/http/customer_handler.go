package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/customer/domain"
	"github.com/ydoro/wishlist/pkg/wishlist/presentation/inputs"
	"github.com/ydoro/wishlist/pkg/wishlist/presentation/outputs"
)

type CunstomerHandler struct {
	createCustomerUC domain.CustomerUC
}

func NewCustomerHandler(r *gin.RouterGroup, uc domain.CustomerUC) {

	handler := &CunstomerHandler{
		createCustomerUC: uc,
	}

	customerRoutes := r.Group("/customers")
	customerRoutes.POST("", handler.CreateCustomer)

}

// CreateCustomer godoc
// @Summary Creates a new customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body inputs.CreateCustomerRequest true "Client data"
// @Success 201 {object} domain.Customer
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/v1/customers [post]
//
// NOTE - Handlers can be moved to a separated layer and vendor agnostic and here we can use some adapters to connect the router with the application
func (h *CunstomerHandler) CreateCustomer(c *gin.Context) {
	var customer inputs.CreateCustomerRequest

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(401, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	id, err := h.createCustomerUC.CreateCustomerWithEmail(c, customer.Email, customer.Name)

	if err != nil {
		fmt.Println("Failed to create customer: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(201, outputs.CreateCustomerResponse{
		ID: id,
	})
}
