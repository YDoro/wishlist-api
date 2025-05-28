package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/usecase"
	mocks "github.com/ydoro/wishlist/mock/domain"
	"go.uber.org/mock/gomock"
)

func TestListCustomerWishlistsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerRepoMock := mocks.NewMockGetCustomerByIDRepository(ctrl)
	wishlistRepoMock := mocks.NewMockWishlistByCustomerIdRepository(ctrl)
	productGetterMock := mocks.NewMockGetProductUseCase(ctrl)

	customer := &domain.Customer{
		ID:    "customer1",
		Name:  "Customer 1",
		Email: "customer1@test.com",
	}

	wishlists := []*domain.Wishlist{
		{
			ID:    "wishlist1",
			Title: "Wishlist 1",
			Items: []string{"product1", "product2"},
		},
		{
			ID:    "wishlist2",
			Title: "Wishlist 2",
			Items: []string{"product3"},
		},
	}

	products := map[string]*domain.Product{
		"product1": {ID: "product1", Name: "Product 1"},
		"product2": {ID: "product2", Name: "Product 2"},
		"product3": {ID: "product3", Name: "Product 3"},
	}

	tests := []struct {
		name              string
		currentCustomerId string
		customerId        string
		mockSetup         func()
		expectedError     error
		expectedLen       int
	}{
		{
			name:              "should return unauthorized when customer ids don't match",
			currentCustomerId: "customer2",
			customerId:        "customer1",
			mockSetup:         func() {},
			expectedError:     e.NewUnauthorizedError(),
			expectedLen:       0,
		},
		{
			name:              "should return error when customer repository fails",
			currentCustomerId: "customer1",
			customerId:        "customer1",
			mockSetup: func() {
				customerRepoMock.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			expectedLen:   0,
		},
		{
			name:              "should return not found when customer doesn't exist",
			currentCustomerId: "customer1",
			customerId:        "customer1",
			mockSetup: func() {
				customerRepoMock.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("customer"),
			expectedLen:   0,
		},
		{
			name:              "should return error when wishlist repository fails",
			currentCustomerId: "customer1",
			customerId:        "customer1",
			mockSetup: func() {
				customerRepoMock.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(customer, nil)

				wishlistRepoMock.EXPECT().
					GetByCustomerId(gomock.Any(), "customer1").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			expectedLen:   0,
		},
		{
			name:              "should return empty list when customer has no wishlists",
			currentCustomerId: "customer1",
			customerId:        "customer1",
			mockSetup: func() {
				customerRepoMock.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(customer, nil)

				wishlistRepoMock.EXPECT().
					GetByCustomerId(gomock.Any(), "customer1").
					Return([]*domain.Wishlist{}, nil)
			},
			expectedError: nil,
			expectedLen:   0,
		},
		{
			name:              "should return error when product getter fails",
			currentCustomerId: "customer1",
			customerId:        "customer1",
			mockSetup: func() {
				customerRepoMock.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(customer, nil)

				wishlistRepoMock.EXPECT().
					GetByCustomerId(gomock.Any(), "customer1").
					Return(wishlists, nil)

				totalProducts := 0
				for _, w := range wishlists {
					totalProducts += len(w.Items)
				}

				productGetterMock.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("product error")).
					MinTimes(1).
					MaxTimes(totalProducts)
			},
			expectedError: errors.New("product error"),
			expectedLen:   0,
		},
		{
			name:              "should return filled wishlists successfully",
			currentCustomerId: "customer1",
			customerId:        "customer1",
			mockSetup: func() {
				customerRepoMock.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(customer, nil)

				wishlistRepoMock.EXPECT().
					GetByCustomerId(gomock.Any(), "customer1").
					Return(wishlists, nil)

				for pid, product := range products {
					productGetterMock.EXPECT().
						Execute(gomock.Any(), pid).
						Return(product, nil).
						Times(1)
				}
			},
			expectedError: nil,
			expectedLen:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			sut := usecase.NewListCustomerWishlistsUseCase(
				customerRepoMock,
				wishlistRepoMock,
				productGetterMock,
			)

			result, err := sut.Execute(context.Background(), tt.currentCustomerId, tt.customerId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLen, len(*result))

				if tt.expectedLen > 0 {
					assert.Equal(t, wishlists[0].ID, (*result)[0].ID)
					assert.Equal(t, customer.ID, (*result)[0].Customer.ID)
					assert.Equal(t, len(wishlists[0].Items), len((*result)[0].Items))
				}
			}
		})
	}
}
