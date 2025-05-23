package domain

import "context"

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
}

type CustomerUC interface {
	CreateCustomerWithEmail(ctx context.Context, email string, name string) (string, error)
}
