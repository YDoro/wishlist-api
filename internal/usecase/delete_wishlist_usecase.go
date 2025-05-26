package usecase

import (
	"context"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type DeleteWishlistUseCase struct {
	CustomerGetter  domain.GetCustomerByIDRepository
	WishlistGetter  domain.WishlistByIdRepository
	WishlistDeleter domain.DeleteWishlistRepository
}

func NewDeleteWishlistUseCase(
	customerGetter domain.GetCustomerByIDRepository,
	wishlistGetter domain.WishlistByIdRepository,
	wishlistDeleter domain.DeleteWishlistRepository,
) *DeleteWishlistUseCase {
	return &DeleteWishlistUseCase{
		CustomerGetter:  customerGetter,
		WishlistGetter:  wishlistGetter,
		WishlistDeleter: wishlistDeleter,
	}
}

func (u *DeleteWishlistUseCase) DeleteWishlist(ctx context.Context, currentCustomerId string, customerId string, wishlistId string) error {
	if currentCustomerId != customerId {
		return e.NewUnauthorizedError()
	}

	customer, err := u.CustomerGetter.GetByID(ctx, customerId)
	if err != nil {
		return err
	}

	if customer == nil {
		return e.NewNotFoundError("customer")
	}

	wishlist, err := u.WishlistGetter.GetById(ctx, wishlistId)
	if err != nil {
		return err
	}

	if wishlist == nil {
		return e.NewNotFoundError("wishlist")
	}

	if wishlist.CustomerId != customerId {
		return e.NewUnauthorizedError()
	}

	return u.WishlistDeleter.DeleteWishlist(ctx, wishlistId)
}
