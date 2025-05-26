package usecase

import (
	"context"
	"time"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type UpdateCustomerUseCase struct {
	Updater     domain.UpdateCustomerRepository
	Getter      domain.GetCustomerByIDRepository
	EmailGetter domain.GetCustomerByEmailRepository
}

func NewUpdateCustomerUseCase(
	updater domain.UpdateCustomerRepository,
	getter domain.GetCustomerByIDRepository,
	emailGetter domain.GetCustomerByEmailRepository,
) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{
		Updater:     updater,
		Getter:      getter,
		EmailGetter: emailGetter,
	}
}

func (u *UpdateCustomerUseCase) UpdateCustomer(ctx context.Context, currentCustomerID string, customerID string, data domain.CustomerEditableFields) (*domain.OutgoingCustomer, error) {
	if currentCustomerID != customerID {
		return nil, e.NewUnauthorizedError()
	}
	customer, err := u.Getter.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if data.Name != "" {
		customer.Name = data.Name
	}

	if data.Email != "" {
		if customer.Email != data.Email {
			existingCustomer, err := u.EmailGetter.GetByEmail(ctx, data.Email)
			if err != nil {
				return nil, err
			}

			if existingCustomer != nil && existingCustomer.ID != customerID {
				return nil, &e.ValidationError{Field: "email", Err: "email already in use"}
			}

			customer.Email = data.Email
		}
	}

	customer.UpdatedAt = time.Now()
	err = u.Updater.Update(ctx, customer)

	if err != nil {
		return nil, err
	}

	return &domain.OutgoingCustomer{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
	}, nil

}
