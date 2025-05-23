package postgresDB

import (
	"context"
	"database/sql"

	"github.com/ydoro/wishlist/internal/customer/domain"
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
	query := `INSERT INTO customers (id, name, email) VALUES ($1, $2, $3)`
	_, err := r.DB.ExecContext(ctx, query, customer.ID, customer.Name, customer.Email)

	return err
}
