package usecase

import (
	"context"
	"fmt"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type ShowWishlistUseCase struct {
	wishlistGetter domain.WishlistByIdRepository
	customerGetter domain.GetCustomerByIDRepository
	productGetter  domain.GetProductUseCase
}

func NewShowWishlistUseCase(
	wishlistGetter domain.WishlistByIdRepository,
	customerGetter domain.GetCustomerByIDRepository,
	productGetter domain.GetProductUseCase,

) *ShowWishlistUseCase {
	return &ShowWishlistUseCase{
		wishlistGetter: wishlistGetter,
		customerGetter: customerGetter,
		productGetter:  productGetter,
	}
}

func (u *ShowWishlistUseCase) ShowWishlist(ctx context.Context, currentCustomerId string, customerId string, wishlistId string) (*domain.FullfilledWishlist, error) {
	if currentCustomerId != customerId {
		return nil, e.NewUnauthorizedError()
	}

	customer, err := u.customerGetter.GetByID(ctx, customerId)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, e.NewNotFoundError("customer")
	}

	wishlist, err := u.wishlistGetter.GetById(ctx, wishlistId)
	if err != nil {
		return nil, err
	}

	if wishlist == nil {
		return nil, e.NewNotFoundError("wishlist")
	}

	if wishlist.CustomerId != customerId {
		return nil, e.NewUnauthorizedError()
	}

	ffwl := &domain.FullfilledWishlist{
		ID: wishlist.ID,
		Customer: &domain.OutgoingCustomer{
			ID:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			CreatedAt: customer.CreatedAt,
		},
		Title: wishlist.Title,
		Items: make([]domain.Product, 0, len(wishlist.Items)),
	}

	if len(wishlist.Items) == 0 {
		return ffwl, nil
	}

	items := u.fetchProductsConcurrently(ctx, wishlist.Items)

	ffwl.Items = items

	return ffwl, nil
}

func (u *ShowWishlistUseCase) fetchProductsConcurrently(ctx context.Context, itemIDs []string) []domain.Product {
	type productResult struct {
		product *domain.Product
		err     error
	}
	results := make(chan productResult, len(itemIDs))
	items := make([]domain.Product, 0, len(itemIDs))

	for _, itemID := range itemIDs {
		go func(id string) {
			select {
			case <-ctx.Done():
				results <- productResult{err: ctx.Err()}
			default:
				product, err := u.productGetter.Execute(ctx, id)
				results <- productResult{product: product, err: err}
			}
		}(itemID)
	}

	for i := 0; i < len(itemIDs); i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("Context cancelled: %v\n", ctx.Err())
			continue
		case result := <-results:
			if result.err != nil {
				fmt.Printf("Error fetching product: %v\n", result.err)
				continue
			}
			if result.product != nil {
				items = append(items, *result.product)
			}
		}
	}

	return items
}
