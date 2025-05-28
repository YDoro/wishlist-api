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
	Customer *OutgoingCustomer
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

type ListUserWishlists interface {
	Execute(ctx context.Context, currentCustomerId string, customerId string) (*[]FullfilledWishlist, error)
}

type DeleteWishlistUseCase interface {
	DeleteWishlist(ctx context.Context, currentCustomerId string, customerId string, wishlistId string) error
}

type UpdateWishListUseCase interface {
	UpdateWishlist(ctx context.Context, currentCustomerId string, wishlist *Wishlist) error
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

type UpdateWishlistRepository interface {
	Update(ctx context.Context, wishlist *Wishlist) error
}

type DeleteWishlistRepository interface {
	DeleteWishlist(ctx context.Context, wishlistId string) error
}
