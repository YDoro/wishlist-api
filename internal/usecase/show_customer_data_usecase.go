package usecase

import (
	"context"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type GetCustomerData struct {
	Getter domain.GetCustomerByIDRepository
}

func NewGetCustomerData(getter domain.GetCustomerByIDRepository) *GetCustomerData {
	return &GetCustomerData{
		Getter: getter,
	}
}

func (g *GetCustomerData) ShowCustomerData(ctx context.Context, currentCustomerId string, id string) (*domain.OutgoingCustomer, error) {
	if currentCustomerId != id {
		return nil, e.NewUnauthorizedError()
	}
	customer, err := g.Getter.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, e.NewNotFoundError("customer")
	}

	out := &domain.OutgoingCustomer{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
	}

	return out, nil
}
