package usecase

import (
	"context"

	"github.com/ydoro/wishlist/internal/domain"
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
	// NOTE - here we can check if some usecase specific authentication logic using the currentCustomerId
	customer, err := g.Getter.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	out := &domain.OutgoingCustomer{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
	}

	return out, nil
}
