package http

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/presentation/errors"
	"github.com/ydoro/wishlist/internal/presentation/inputs"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

type CunstomerHandler struct {
	createCustomerUC domain.CreateCustomerUC
	showcustomerUC   domain.ShowCustomerDataUC
}

func NewCustomerHandler(
	r *gin.RouterGroup,
	uc domain.CreateCustomerUC,
	auth gin.HandlerFunc,
	showCustomeruc domain.ShowCustomerDataUC,
) {
	handler := &CunstomerHandler{
		createCustomerUC: uc,
		showcustomerUC:   showCustomeruc,
	}

	customerRoutes := r.Group("/customers")
	customerRoutes.GET("/:id", auth, handler.ShowCustomerData)
	customerRoutes.POST("/", handler.CreateCustomer)

}

// CreateCustomer godoc
// @Summary Creates a new customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body inputs.CreateCustomerRequest true "Client data"
// @Success 201 {object} domain.Customer
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers [post]
//
// NOTE - Handlers can be moved to a separated layer and vendor agnostic and here we can use some adapters to connect the router with the application
func (h *CunstomerHandler) CreateCustomer(c *gin.Context) {
	var customer inputs.CreateCustomerRequest

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	id, err := h.createCustomerUC.CreateCustomerWithEmail(c, domain.IncommingCustomer(customer))

	if err != nil {
		if e.IsValidationError(err) {
			c.JSON(400, outputs.ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		fmt.Printf("\nFailed to create customer: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(201, outputs.CreateCustomerResponse{
		ID: id,
	})
}

// GetCustomerData godoc
// @Summary Gets given customer data by ID
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.OutgoingCustomer
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers [get]
func (h *CunstomerHandler) ShowCustomerData(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	var customer domain.OutgoingCustomer

	str, _ := c.Get("currentCustomer")
	json.Unmarshal([]byte(str.(string)), &customer)

	res, err := h.showcustomerUC.ShowCustomerData(c, customer.ID, id)

	if err != nil {
		if e.IsUnauthorizedError(err) {
			c.JSON(401, outputs.ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to get customer"})
		return
	}
	c.JSON(200, res)
}
