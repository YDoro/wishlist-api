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

type CunstomerHandler struct {
	createCustomerUC domain.CreateCustomerUC
	showcustomerUC   domain.ShowCustomerDataUC
	updatecustomerUC domain.UpdateCustomerUC
	deleteCustomerUC domain.DeleteCustomerUC
}

func NewCustomerHandler(
	r *gin.RouterGroup,
	uc domain.CreateCustomerUC,
	auth gin.HandlerFunc,
	showCustomeruc domain.ShowCustomerDataUC,
	updateCustomerUC domain.UpdateCustomerUC,
	deleteCustomerUC domain.DeleteCustomerUC,
) {
	handler := &CunstomerHandler{
		createCustomerUC: uc,
		showcustomerUC:   showCustomeruc,
		updatecustomerUC: updateCustomerUC,
		deleteCustomerUC: deleteCustomerUC,
	}

	customerRoutes := r.Group("/customers")
	customerRoutes.GET("/:id", auth, handler.ShowCustomerData)
	customerRoutes.PATCH("/:id", auth, handler.UpdateCustomer)
	customerRoutes.DELETE("/:id", auth, handler.DeleteCustomer)
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
// @Param id path string true "Customer ID"
// @Success 200 {object} domain.OutgoingCustomer
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{id} [get]
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

// UpdateCustomer godoc
// @Summary updates the given customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body domain.CustomerEditableFields true "Data to update, either name or email or both"
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} domain.OutgoingCustomer
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{id} [patch]
func (h *CunstomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	var currentCustomer domain.OutgoingCustomer

	str, _ := c.Get("currentCustomer")
	json.Unmarshal([]byte(str.(string)), &currentCustomer)

	var data domain.CustomerEditableFields
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	res, err := h.updatecustomerUC.UpdateCustomer(c, currentCustomer.ID, id, data)

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

// DeleteCustomer godoc
// @Summary Deletes the given customer
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 204
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 404 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{id} [delete]
func (h *CunstomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	var currentCustomer domain.OutgoingCustomer

	str, _ := c.Get("currentCustomer")
	json.Unmarshal([]byte(str.(string)), &currentCustomer)

	err := h.deleteCustomerUC.DeleteCustomer(c, currentCustomer.ID, id)

	if err != nil {
		switch {
		case e.IsUnauthorizedError(err):
			c.JSON(401, outputs.ErrorResponse{
				Message: err.Error(),
			})
		case e.IsNotFoundError(err):
			c.JSON(404, outputs.ErrorResponse{
				Message: err.Error(),
			})
		default:
			c.JSON(500, outputs.ErrorResponse{
				Message: "Failed to delete customer",
			})
		}
		return
	}

	c.Status(204)
}
