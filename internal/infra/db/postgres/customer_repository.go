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

func (r *customerRepo) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM customers WHERE email = $1 AND deleted_at IS NULL`
	row := r.DB.QueryRowContext(ctx, query, email)

	customer := &domain.Customer{}
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Password, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return customer, nil
}

func (r *customerRepo) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM customers WHERE id = $1 AND deleted_at IS NULL`
	row := r.DB.QueryRowContext(ctx, query, id)

	customer := &domain.Customer{}
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Password, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return customer, nil
}

func (r *customerRepo) Update(ctx context.Context, customer *domain.Customer) error {
	query := `UPDATE customers 
		SET name = $1, 
			email = $2, 
			updated_at = $3, 
			deleted_at = $4 
		WHERE id = $5`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		customer.Name,
		customer.Email,
		customer.UpdatedAt,
		customer.DeletedAt,
		customer.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
