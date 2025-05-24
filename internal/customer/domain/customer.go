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

type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
}

type CustomerUC interface {
	CreateCustomerWithEmail(ctx context.Context, data IncommingCustomer) (string, error)
}
