package usecase

import (
	"context"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type CreateWishlistUseCase struct {
	CustomerGetter domain.GetCustomerByIDRepository
	Getter         domain.WishlistByTitleRepository
	Creator        domain.WishlistCreationRepository
	IdMaker        domain.IDGenerator
}

func NewCreateWishlistUseCase(
	getter domain.WishlistByTitleRepository,
	creator domain.WishlistCreationRepository,
	customerGetter domain.GetCustomerByIDRepository,
) *CreateWishlistUseCase {
	return &CreateWishlistUseCase{
		Getter:         getter,
		Creator:        creator,
		CustomerGetter: customerGetter,
	}
}
func (u *CreateWishlistUseCase) CreateWishlist(ctx context.Context, currentCustomerId string, customerId string, title string) (string, error) {
	if currentCustomerId != customerId {
		return "", e.NewUnauthorizedError()
	}

	customer, err := u.CustomerGetter.GetByID(ctx, customerId)
	if err != nil {
		return "", err
	}

	if customer == nil {
		return "", e.NewNotFoundError("customer")
	}

	wishlist, err := u.Getter.GetByTitle(ctx, customerId, title)
	if err != nil {
		return "", err
	}

	if wishlist != nil {
		return "", &e.ValidationError{
			Field: "title",
			Err:   "already in use",
		}
	}

	newId, err := u.IdMaker.Generate()
	if err != nil {
		return "", err

	}

	newWishlist := &domain.Wishlist{
		ID:         newId,
		CustomerId: customerId,
		Title:      title,
		Items:      []string{},
	}

	err = u.Creator.Create(ctx, newWishlist)
	if err != nil {
		return "", err
	}

	return newId, nil

}
