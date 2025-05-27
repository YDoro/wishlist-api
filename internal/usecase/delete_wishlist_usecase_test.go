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

func TestDeleteWishlistUseCase_DeleteWishlist(t *testing.T) {
	tests := []struct {
		name              string
		currentCustomerId string
		customerId        string
		wishlistId        string
		setupMocks        func(*mocks.MockGetCustomerByIDRepository, *mocks.MockWishlistByIdRepository, *mocks.MockDeleteWishlistRepository)
		expectedError     error
	}{
		{
			name:              "successful wishlist deletion",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil).
					Times(1)

				wishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(&domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "customer_123",
						Title:      "Birthday Wishlist",
						Items:      []string{},
					}, nil).
					Times(1)

				wishlistDeleter.EXPECT().
					DeleteWishlist(gomock.Any(), "wishlist_123").
					Return(nil).
					Times(1)
			},
			expectedError: nil,
		},
		{
			name:              "unauthorized deletion attempt",
			currentCustomerId: "different_customer",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
			},
			expectedError: e.NewUnauthorizedError(),
		},
		{
			name:              "customer not found",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("customer"),
		},
		{
			name:              "customer getter error",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "wishlist not found",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("wishlist"),
		},
		{
			name:              "wishlist getter error",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "wishlist belongs to different customer",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(&domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "different_customer",
						Title:      "Birthday Wishlist",
						Items:      []string{},
					}, nil)
			},
			expectedError: e.NewUnauthorizedError(),
		},
		{
			name:              "deletion error",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByIdRepository,
				wishlistDeleter *mocks.MockDeleteWishlistRepository) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(&domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "customer_123",
						Title:      "Birthday Wishlist",
						Items:      []string{},
					}, nil)

				wishlistDeleter.EXPECT().
					DeleteWishlist(gomock.Any(), "wishlist_123").
					Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
			wishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
			wishlistDeleter := mocks.NewMockDeleteWishlistRepository(ctrl)

			tt.setupMocks(customerGetter, wishlistGetter, wishlistDeleter)

			uc := usecase.NewDeleteWishlistUseCase(
				customerGetter,
				wishlistGetter,
				wishlistDeleter,
			)

			err := uc.DeleteWishlist(context.Background(), tt.currentCustomerId, tt.customerId, tt.wishlistId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
