package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type UpdateWishListUseCase struct {
	customerRepository domain.GetCustomerByIDRepository
	getterRepository   domain.WishlistByIdRepository
	updateRepository   domain.UpdateWishlistRepository
	productGetter      domain.GetProductUseCase
}

func NewUpdateWishListUseCase(
	customerRepository domain.GetCustomerByIDRepository,
	getterRepository domain.WishlistByIdRepository,
	updateRepository domain.UpdateWishlistRepository,
	productGetter domain.GetProductUseCase,
) *UpdateWishListUseCase {
	return &UpdateWishListUseCase{
		customerRepository: customerRepository,
		getterRepository:   getterRepository,
		updateRepository:   updateRepository,
		productGetter:      productGetter,
	}
}

func (u *UpdateWishListUseCase) UpdateWishlist(ctx context.Context, currentCustomerId string, wishlist *domain.Wishlist) error {
	if currentCustomerId != wishlist.CustomerId {
		return e.NewUnauthorizedError()
	}

	customer, err := u.customerRepository.GetByID(ctx, wishlist.CustomerId)
	if err != nil {
		return err
	}

	if customer == nil {
		return e.NewNotFoundError("customer")
	}

	dbWishlist, err := u.getterRepository.GetById(ctx, wishlist.ID)
	if err != nil {
		return err
	}

	if dbWishlist == nil {
		return e.NewNotFoundError("wishlist")
	}

	if dbWishlist.CustomerId != currentCustomerId {
		return e.NewUnauthorizedError()
	}

	if wishlist.Title == dbWishlist.Title && slices.Equal(wishlist.Items, dbWishlist.Items) {
		return nil
	}

	if wishlist.Title != "" {
		dbWishlist.Title = wishlist.Title
	}

	if wishlist.Items != nil {
		var newProducts []string
		for _, productId := range wishlist.Items {
			product, err := u.productGetter.Execute(ctx, productId)

			if err != nil {
				fmt.Printf("Error fetching product: %s %v\n", productId, err)
				return err
			}
			if product == nil {
				return e.NewNotFoundError(fmt.Sprintf("product_%s", productId))
			}

			newProducts = append(newProducts, productId)
		}

		dbWishlist.Items = newProducts
	}

	return u.updateRepository.Update(ctx, dbWishlist)
}
