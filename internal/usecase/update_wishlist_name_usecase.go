package usecase

import (
	"context"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type UpdateWishlistNameUseCase struct {
	CustomerGetter  domain.GetCustomerByIDRepository
	WishlistGetter  domain.WishlistByIdRepository
	WishlistUpdater domain.UpdateWishlistNameRepository
}

func NewUpdateWishlistNameUseCase(
	customerGetter domain.GetCustomerByIDRepository,
	wishlistGetter domain.WishlistByIdRepository,
	wishlistUpdater domain.UpdateWishlistNameRepository,
) *UpdateWishlistNameUseCase {
	return &UpdateWishlistNameUseCase{
		CustomerGetter:  customerGetter,
		WishlistGetter:  wishlistGetter,
		WishlistUpdater: wishlistUpdater,
	}
}

func (u *UpdateWishlistNameUseCase) Rename(ctx context.Context, currentCustomerId string, customerId string, wishlistId string, title string) (string, error) {
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

	wishlist, err := u.WishlistGetter.GetById(ctx, wishlistId)
	if err != nil {
		return "", err
	}

	if wishlist == nil {
		return "", e.NewNotFoundError("wishlist")
	}

	if wishlist.CustomerId != customerId {
		return "", e.NewUnauthorizedError()
	}

	wishlist.Title = title

	err = u.WishlistUpdater.UpdateWishlistName(ctx, wishlist)
	if err != nil {
		return "", err
	}

	return wishlist.ID, nil
}
