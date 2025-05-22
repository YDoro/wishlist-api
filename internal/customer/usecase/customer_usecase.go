package usecase

import (
	"context"

	"github.com/ydoro/wishlist/internal/customer/domain"
)

type CustomerUseCase struct {
	repo domain.CustomerRepository
}

func NewCustomerUseCase(customerRepository domain.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		repo: customerRepository,
	}
}

func (uc *CustomerUseCase) CreateCustomerWithEmail(ctx context.Context, email string, name string) (string, error) {
	// TODO - do not expose db id
	return uc.repo.Create(ctx, &domain.Customer{
		Name:  name,
		Email: email,
		ID:    nil,
	})
}
