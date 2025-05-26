package usecase

import (
	"context"
	"time"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/presentation/errors"
)

type DeleteCustomerUseCase struct {
	Updater domain.UpdateCustomerRepository
	Getter  domain.GetCustomerByIDRepository
}

func NewDeleteCustomerUseCase(
	updater domain.UpdateCustomerRepository,
	getter domain.GetCustomerByIDRepository,
) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{
		Updater: updater,
		Getter:  getter,
	}
}

func (u *DeleteCustomerUseCase) DeleteCustomer(ctx context.Context, currentCustomerID string, customerID string) error {
	if currentCustomerID != customerID {
		return e.NewUnauthorizedError()
	}

	customer, err := u.Getter.GetByID(ctx, customerID)
	if err != nil {
		return err
	}

	if customer == nil {
		return e.NewNotFoundError("customer")
	}

	customer.DeletedAt = time.Now()
	err = u.Updater.Update(ctx, customer)
	if err != nil {
		return err
	}

	return nil
}
