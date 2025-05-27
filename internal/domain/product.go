//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/product_mock.go -package=mocks -source ./product.go

package domain

import "context"

type Rating struct {
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
	Rating      *Rating  `json:"rating"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	DeletedAt   string   `json:"deleted_at,omitempty"`
}

type GetProductUseCase interface {
	Execute(ctx context.Context, productID string) (*Product, error)
}

// both external service and database can implement these interfaces
type GetProductRepository interface {
	GetByID(ctx context.Context, productID string) (*Product, error)
}

type ListProductsRepository interface {
	List(ctx context.Context, count int, offset int) ([]Product, error)
}

type UpsertProductRepository interface {
	Upsert(ctx context.Context, product Product) error
}

type DeleteProductRepository interface {
	Delete(ctx context.Context, productID string) error
}
