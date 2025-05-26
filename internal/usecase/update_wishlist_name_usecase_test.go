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

func TestUpdateWishlistNameUseCase_Rename(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
	mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
	mockWishlistUpdater := mocks.NewMockUpdateWishlistNameRepository(ctrl)

	tests := []struct {
		name              string
		currentCustomerId string
		customerId        string
		wishlistId        string
		title             string
		setupMocks        func()
		expectedId        string
		expectedError     error
	}{
		{
			name:              "successful wishlist update",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks: func() {
				mockCustomerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				mockWishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(&domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "customer_123",
						Title:      "Birthday Wishlist",
						Items:      []string{},
					}, nil)

				mockWishlistUpdater.EXPECT().
					UpdateWishlistName(gomock.Any(), &domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "customer_123",
						Title:      "Updated Birthday Wishlist",
						Items:      []string{},
					}).
					Return(nil)
			},
			expectedId:    "wishlist_123",
			expectedError: nil,
		},
		{
			name:              "unauthorized update attempt",
			currentCustomerId: "different_customer",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks:        func() {},
			expectedId:        "",
			expectedError:     e.NewUnauthorizedError(),
		},
		{
			name:              "customer not found",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks: func() {
				mockCustomerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, nil)
			},
			expectedId:    "",
			expectedError: e.NewNotFoundError("customer"),
		},
		{
			name:              "customer getter error",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks: func() {
				mockCustomerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, errors.New("database error"))
			},
			expectedId:    "",
			expectedError: errors.New("database error"),
		},
		{
			name:              "wishlist not found",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks: func() {
				mockCustomerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				mockWishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(nil, nil)
			},
			expectedId:    "",
			expectedError: e.NewNotFoundError("wishlist"),
		},
		{
			name:              "wishlist getter error",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks: func() {
				mockCustomerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				mockWishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(nil, errors.New("database error"))
			},
			expectedId:    "",
			expectedError: errors.New("database error"),
		},
		{
			name:              "wishlist belongs to different customer",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks: func() {
				mockCustomerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				mockWishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(&domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "different_customer",
						Title:      "Birthday Wishlist",
						Items:      []string{},
					}, nil)
			},
			expectedId:    "",
			expectedError: e.NewUnauthorizedError(),
		},
		{
			name:              "update error",
			currentCustomerId: "customer_123",
			customerId:        "customer_123",
			wishlistId:        "wishlist_123",
			title:             "Updated Birthday Wishlist",
			setupMocks: func() {
				mockCustomerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				mockWishlistGetter.EXPECT().
					GetById(gomock.Any(), "wishlist_123").
					Return(&domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "customer_123",
						Title:      "Birthday Wishlist",
						Items:      []string{},
					}, nil)

				mockWishlistUpdater.EXPECT().
					UpdateWishlistName(gomock.Any(), &domain.Wishlist{
						ID:         "wishlist_123",
						CustomerId: "customer_123",
						Title:      "Updated Birthday Wishlist",
						Items:      []string{},
					}).
					Return(errors.New("database error"))
			},
			expectedId:    "",
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			uc := usecase.NewUpdateWishlistNameUseCase(
				mockCustomerGetter,
				mockWishlistGetter,
				mockWishlistUpdater,
			)

			id, err := uc.Rename(context.Background(), tt.currentCustomerId, tt.customerId, tt.wishlistId, tt.title)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedId, id)
		})
	}
}
