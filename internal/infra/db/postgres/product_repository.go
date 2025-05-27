package postgresDB

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/ydoro/wishlist/internal/domain"
)

type productRepo struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepo {
	return &productRepo{
		DB: db,
	}
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (
			id, name, price, description, images, rating, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)`

	var ratingJSON []byte
	var err error
	if product.Rating != nil {
		ratingJSON, err = json.Marshal(product.Rating)
		if err != nil {
			return err
		}
	}

	_, err = r.DB.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Price,
		product.Description,
		pq.Array(product.Images),
		ratingJSON,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *productRepo) GetByID(ctx context.Context, productID string) (*domain.Product, error) {
	query := `SELECT id, name, price, description, images, rating, created_at, updated_at FROM products WHERE id = $1 AND deleted_at IS NULL`
	row := r.DB.QueryRowContext(ctx, query, productID)

	product := &domain.Product{}
	var ratingJSON []byte
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Description,
		pq.Array(&product.Images),
		&ratingJSON,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if ratingJSON != nil {
		var rating domain.Rating
		if err := json.Unmarshal(ratingJSON, &rating); err != nil {
			return nil, err
		}
		product.Rating = &rating
	}

	return product, nil
}
