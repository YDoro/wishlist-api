package usecase

import (
	"context"
	"slices"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type AddProductToWishListUC struct {
	customerGetter  domain.GetCustomerByIDRepository
	wishlistGetter  domain.WishlistByIdRepository
	wishlistUpdater domain.UpdateWishlistRepository
	productGetter   domain.GetProductUseCase
}

func NewAddProductToWishListUC(
	customerGetter domain.GetCustomerByIDRepository,
	wishlistGetter domain.WishlistByIdRepository,
	wishlistUpdater domain.UpdateWishlistRepository,
	productGetter domain.GetProductUseCase,
) *AddProductToWishListUC {
	return &AddProductToWishListUC{
		customerGetter:  customerGetter,
		wishlistGetter:  wishlistGetter,
		wishlistUpdater: wishlistUpdater,
		productGetter:   productGetter,
	}
}

func (uc *AddProductToWishListUC) AddProduct(ctx context.Context, currentCustomerId string, customerId string, wishlistId string, productId string) error {
	if currentCustomerId != customerId {
		return e.NewUnauthorizedError()
	}

	customer, err := uc.customerGetter.GetByID(ctx, customerId)
	if err != nil {
		return err
	}
	if customer == nil {
		return e.NewNotFoundError("customer")
	}

	wishlist, err := uc.wishlistGetter.GetById(ctx, wishlistId)
	if err != nil {
		return err
	}
	if wishlist == nil {
		return e.NewNotFoundError("wishlist")
	}

	if wishlist.CustomerId != customerId {
		return e.NewUnauthorizedError()
	}

	if slices.Contains(wishlist.Items, productId) {
		return nil
	}

	product, err := uc.productGetter.Execute(ctx, productId)

	if err != nil {
		return err
	}
	if product == nil {
		return e.NewNotFoundError("product")
	}

	wishlist.Items = append(wishlist.Items, productId)
	// TODO - add an updated_at
	if err := uc.wishlistUpdater.Update(ctx, wishlist); err != nil {
		return err
	}
	return nil
}
