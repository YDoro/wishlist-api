//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/customer_mock.go -package=mocks . CreateCustomerUC,ShowCustomerDataUC,CustomerCreationRepository,GetCustomerByEmailRepository,GetCustomerByIDRepository,UpdateCustomerUC,UpdateCustomerRepository,DeleteCustomerUC

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

type CustomerEditableFields struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	// NOTE - Password is intentionally omitted here to prevent accidental updates
	// TODO - Password changes should be handled separately
}

type OutgoingCustomer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Usecases
type CreateCustomerUC interface {
	CreateCustomerWithEmail(ctx context.Context, data IncommingCustomer) (string, error)
}

type ShowCustomerDataUC interface {
	ShowCustomerData(ctx context.Context, currentCustomerid string, id string) (*OutgoingCustomer, error)
}

type UpdateCustomerUC interface {
	UpdateCustomer(ctx context.Context, currentCustomerID string, customerID string, data CustomerEditableFields) (*OutgoingCustomer, error)
}

type DeleteCustomerUC interface {
	DeleteCustomer(ctx context.Context, currentCustomerID string, customerID string) error
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

type UpdateCustomerRepository interface {
	Update(ctx context.Context, customer *Customer) error
}
