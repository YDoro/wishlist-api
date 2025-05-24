package postgresDB

import (
	"context"
	"database/sql"

	"github.com/ydoro/wishlist/internal/domain"
)

type customerRepo struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *customerRepo {
	return &customerRepo{
		DB: db,
	}
}

func (r *customerRepo) Create(ctx context.Context, customer *domain.Customer) error {
	query := `INSERT INTO customers (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.ExecContext(ctx, query, customer.ID, customer.Name, customer.Email, customer.Password, customer.CreatedAt, customer.UpdatedAt)

	return err
}
