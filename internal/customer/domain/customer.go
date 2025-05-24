//go:generate mockgen --build_flags=--mod=mod -destination=../../../mock/domain/customer_mock.go -package=mocks . CreateCustomerUC,CustomerCreationRepository,GetCustomerByEmailRepository,GetCustomerByIDRepository

package domain

import (
	"context"
	"time"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type IncommingCustomer struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Usecases

type CreateCustomerUC interface {
	CreateCustomerWithEmail(ctx context.Context, data IncommingCustomer) (string, error)
}

// repositories

type CustomerCreationRepository interface {
	Create(ctx context.Context, customer *Customer) error
}

// NOTE - this getter interfaces couldbe merged into a single one with some abstraction
type GetCustomerByEmailRepository interface {
	GetByEmail(ctx context.Context, email string) (*Customer, error)
}

type GetCustomerByIDRepository interface {
	GetByID(ctx context.Context, id string) (*Customer, error)
}
