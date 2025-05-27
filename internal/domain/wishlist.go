//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/wishlist_mock.go -package=mocks -source ./wishlist.go
package domain

import "context"

type Wishlist struct {
	ID         string   `json:"id"`
	CustomerId string   `json:"customer_id"`
	Title      string   `json:"title"`
	Items      []string `json:"items"`
}

type FullfilledWishlist struct {
	ID       string
	Customer *Customer
	Title    string
	Items    []Product
}

// Usecases

type CreateWishlistUseCase interface {
	CreateWishlist(ctx context.Context, currentCustomerId string, customerId string, title string) (string, error)
}

type ShowWishlistUseCase interface {
	ShowWishlist(ctx context.Context, currentCustomerId string, customerId string, wishlistId string) (*FullfilledWishlist, error)
}

type UpdateWishlistNameUseCase interface {
	RenameWishlist(ctx context.Context, currentCustomerId string, customerId string, wishlistId string, title string) (string, error)
}

type DeleteWishlistUseCase interface {
	DeleteWishlist(ctx context.Context, currentCustomerId string, customerId string, wishlistId string) error
}

type AddProductToWishlistUseCase interface {
	AddProduct(ctx context.Context, currentCustomerId string, customerId string, wishlistId string, productId string) error
}

type RemoveProductToWishlistUseCase interface {
	RemoveProduct(ctx context.Context, currentCustomerId string, customerId string, wishlistId string, productId string) (*FullfilledWishlist, error)
}

// Repositories

type WishlistCreationRepository interface {
	Create(ctx context.Context, wishlist *Wishlist) error
}

type WishlistByIdRepository interface {
	GetById(ctx context.Context, wishlistId string) (*Wishlist, error)
}

type WishlistByTitleRepository interface {
	GetByTitle(ctx context.Context, customerId string, title string) (*Wishlist, error)
}

type WishlistByCustomerIdRepository interface {
	GetByCustomerId(ctx context.Context, customerId string) ([]*Wishlist, error)
}

type UpdateWishlistNameRepository interface {
	UpdateWishlistName(ctx context.Context, wishlist *Wishlist) error
}

type DeleteWishlistRepository interface {
	DeleteWishlist(ctx context.Context, wishlistId string) error
}
type AddItemToWishlistRepository interface {
	AddItemToWishlist(ctx context.Context, wishlistId string, itemId string) error
}

type RemoveItemFromWishlistRepository interface {
	RemoveItemFromWishlist(ctx context.Context, wishlistId string, itemId string) error
}
