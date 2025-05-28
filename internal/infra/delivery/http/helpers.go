package http

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/presentation/outputs"
)

func HandleError(c *gin.Context, err error) {
	if err != nil {
		if e.IsValidationError(err) {
			c.JSON(400, outputs.ErrorResponse{
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
		if e.IsNotFoundError(err) {
			c.JSON(404, outputs.ErrorResponse{
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
}

func GetCustomerFromContext(c *gin.Context) *domain.OutgoingCustomer {
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
