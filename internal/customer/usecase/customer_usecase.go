package usecase

import (
	"context"
	"errors"

	"github.com/ydoro/wishlist/internal/customer/domain"
)

type CustomerUseCase struct {
	repo  domain.CustomerRepository
	idGen domain.IDGenerator
}

func NewCustomerUseCase(customerRepository domain.CustomerRepository, idGen domain.IDGenerator) *CustomerUseCase {
	return &CustomerUseCase{
		repo:  customerRepository,
		idGen: idGen,
	}
}

func (uc *CustomerUseCase) CreateCustomerWithEmail(ctx context.Context, email string, name string) (string, error) {
	id, err := uc.idGen.Generate()

	if err != nil {
		return "", errors.Join(err, errors.New("failed to generate customer ID"))
	}

	return id, uc.repo.Create(ctx, &domain.Customer{
		Name:  name,
		Email: email,
		ID:    id,
	})
}
