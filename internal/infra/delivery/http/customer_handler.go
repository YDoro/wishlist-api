package http

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
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
) *gin.RouterGroup {
	handler := &CunstomerHandler{
		createCustomerUC: uc,
		showcustomerUC:   showCustomeruc,
		updatecustomerUC: updateCustomerUC,
		deleteCustomerUC: deleteCustomerUC,
	}

	customerRoutes := r.Group("/customers")
	customerRoutes.POST("/", handler.CreateCustomer)
	customerRoutes.GET("/:customerId", auth, handler.ShowCustomerData)
	customerRoutes.PATCH("/:customerId", auth, handler.UpdateCustomer)
	customerRoutes.DELETE("/:customerId", auth, handler.DeleteCustomer)

	return customerRoutes

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
		HandleError(c, err)
		return
	}

	c.JSON(201, outputs.CreateCustomerResponse{
		ID: id,
	})
	return
}

// GetCustomerData godoc
// @Summary Gets given customer data by ID
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param customerId path string true "Customer ID"
// @Success 200 {object} domain.OutgoingCustomer
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId} [get]
func (h *CunstomerHandler) ShowCustomerData(c *gin.Context) {
	h.ensureParams(c)
	currentCustomer := GetCustomerFromContext(c)
	res, err := h.showcustomerUC.ShowCustomerData(c, currentCustomer.ID, c.Param("customerId"))

	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, res)
	return
}

// UpdateCustomer godoc
// @Summary updates the given customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body domain.CustomerEditableFields true "Data to update, either name or email or both"
// @Security BearerAuth
// @Param customerId path string true "Customer ID"
// @Success 200 {object} domain.OutgoingCustomer
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId} [patch]
func (h *CunstomerHandler) UpdateCustomer(c *gin.Context) {
	h.ensureParams(c)
	currentCustomer := GetCustomerFromContext(c)

	var data domain.CustomerEditableFields
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		return
	}

	res, err := h.updatecustomerUC.UpdateCustomer(c, currentCustomer.ID, c.Param("customerId"), data)

	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, res)
	return
}

// DeleteCustomer godoc
// @Summary Deletes the given customer
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param customerId path string true "Customer ID"
// @Success 204
// @Failure 400 {object} outputs.ErrorResponse
// @Failure 401 {object} outputs.ErrorResponse
// @Failure 404 {object} outputs.ErrorResponse
// @Failure 500 {object} outputs.ErrorResponse
// @Router /api/customers/{customerId} [delete]
func (h *CunstomerHandler) DeleteCustomer(c *gin.Context) {
	h.ensureParams(c)

	var currentCustomer domain.OutgoingCustomer

	str, _ := c.Get("currentCustomer")
	json.Unmarshal([]byte(str.(string)), &currentCustomer)

	err := h.deleteCustomerUC.DeleteCustomer(c, currentCustomer.ID, c.Param("customerId"))

	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(204)
	return
}

func (h CunstomerHandler) ensureParams(c *gin.Context) {
	cid := c.Param("customerId")

	if cid == "" {
		c.JSON(400, outputs.ErrorResponse{
			Message: "Invalid input"})
		c.Next()
		return
	}
}
