package postgresDB

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *customerRepo) Create(ctx context.Context, customer *domain.Customer) (string, error) {
	query := `INSERT INTO customers (id, name, email) VALUES ($1, $2, $3)`

	var id int64
	err := r.DB.QueryRowContext(ctx, query, customer.ID, customer.Name, customer.Email).Scan(&id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}
