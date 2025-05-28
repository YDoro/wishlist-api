package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type createCustomerUseCase struct {
	repo           domain.CustomerCreationRepository
	idGen          domain.IDGenerator
	pwdHasher      domain.Hasher
	customerGetter domain.GetCustomerByEmailRepository
}

func NewCreateCustomerUseCase(customerRepository domain.CustomerCreationRepository, idGen domain.IDGenerator, hasher domain.Hasher, customerGetter domain.GetCustomerByEmailRepository) *createCustomerUseCase {
	return &createCustomerUseCase{
		repo:           customerRepository,
		idGen:          idGen,
		pwdHasher:      hasher,
		customerGetter: customerGetter,
	}
}

func (uc *createCustomerUseCase) CreateCustomerWithEmail(ctx context.Context, data domain.IncommingCustomer) (string, error) {
	if data.Password == "" {
		return "", e.NewRequiredFieldError("password")
	}

	c, err := uc.customerGetter.GetByEmail(ctx, data.Email)

	if c != nil {
		return "", &e.ValidationError{Field: "email", Err: "email already in use"}
	}

	// here we can add some validation for the password strength and etc
	pwd, err := uc.pwdHasher.Hash(data.Password)
	if err != nil {
		return "", errors.Join(err, errors.New("failed to hash password"))
	}

	id, err := uc.idGen.Generate()
	if err != nil {
		return "", errors.Join(err, errors.New("failed to generate customer ID"))
	}

	return id, uc.repo.Create(ctx, &domain.Customer{
		Name:      data.Name,
		Email:     data.Email,
		Password:  pwd,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        id,
	})
}
