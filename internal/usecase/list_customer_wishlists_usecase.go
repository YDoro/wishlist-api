package usecase

import (
	"context"
	"sync"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type listCustomerWishlistsUseCase struct {
	customerRepo  domain.GetCustomerByIDRepository
	wishlistRepo  domain.WishlistByCustomerIdRepository
	productGetter domain.GetProductUseCase
}

func NewListCustomerWishlistsUseCase(
	customerRepo domain.GetCustomerByIDRepository,
	wishlistRepo domain.WishlistByCustomerIdRepository,
	productGetter domain.GetProductUseCase,
) *listCustomerWishlistsUseCase {
	return &listCustomerWishlistsUseCase{
		customerRepo:  customerRepo,
		wishlistRepo:  wishlistRepo,
		productGetter: productGetter,
	}
}

func (u *listCustomerWishlistsUseCase) Execute(ctx context.Context, currentCustomerId string, customerId string) ([]*domain.FullfilledWishlist, error) {
	if currentCustomerId != customerId {
		return nil, e.NewUnauthorizedError()
	}

	customer, err := u.customerRepo.GetByID(ctx, customerId)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, e.NewNotFoundError("customer")
	}

	wishlists, err := u.wishlistRepo.GetByCustomerId(ctx, customerId)
	if err != nil {
		return nil, err
	}

	return u.fillCustomerWishlists(ctx, wishlists, customer)
}

func (u *listCustomerWishlistsUseCase) fillCustomerWishlists(ctx context.Context, wishlists []*domain.Wishlist, customer *domain.Customer) ([]*domain.FullfilledWishlist, error) {
	if len(wishlists) == 0 {
		return make([]*domain.FullfilledWishlist, 0), nil
	}

	filledLists := make([]*domain.FullfilledWishlist, len(wishlists))
	var wg sync.WaitGroup
	errChan := make(chan error, len(wishlists))

	for i, wishlist := range wishlists {
		wg.Add(1)
		go func(i int, wishlist *domain.Wishlist) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				filledList, err := u.fillWishlistWithProducts(ctx, wishlist, customer)
				if err != nil {
					select {
					case errChan <- err:
					default:
					}
					return
				}

				filledLists[i] = filledList
			}
		}(i, wishlist)
	}

	wg.Wait()
	close(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	default:
	}

	return filledLists, nil
}

func (u *listCustomerWishlistsUseCase) fillWishlistWithProducts(ctx context.Context, wishlist *domain.Wishlist, customer *domain.Customer) (*domain.FullfilledWishlist, error) {
	filledList := &domain.FullfilledWishlist{
		ID: wishlist.ID,
		Customer: &domain.OutgoingCustomer{
			ID:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			CreatedAt: customer.CreatedAt,
		},
		Title: wishlist.Title,
		Items: make([]domain.Product, len(wishlist.Items)),
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(wishlist.Items))

	for i, item := range wishlist.Items {
		wg.Add(1)
		go func(i int, productID string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				product, err := u.productGetter.Execute(ctx, productID)
				if err != nil {
					errChan <- err
					return
				}

				filledList.Items[i] = *product
			}
		}(i, item)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return nil, err
	}

	return filledList, nil
}
