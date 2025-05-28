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

func TestUpdateWishListUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
	mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
	mockProductGetter := mocks.NewMockGetProductUseCase(ctrl)
	mockWishlistUpdater := mocks.NewMockUpdateWishlistRepository(ctrl)

	tests := []struct {
		name              string
		currentCustomerID string
		wishlistTitle     string
		customerID        string
		wishlistID        string
		products          []string
		setupMocks        func()
		expectedError     error
	}{
		{
			name:              "should return unauthorized when current customer is different",
			currentCustomerID: "customer1",
			customerID:        "customer2",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
			},
			expectedError: e.NewUnauthorizedError(),
		},
		{
			name:              "should retunn error on customer get error",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "should return error on customer not found",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("customer"),
		},
		{
			name:              "should return error on wishlist get error",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "should return error on wishlist not found",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("wishlist"),
		},
		{
			name:              "should return eror on other user wishlist",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{ID: "wishlist1", CustomerId: "customer2"}, nil)
			},
			expectedError: e.NewUnauthorizedError(),
		},
		{
			name:              "should succeed if there is nothing to update",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			wishlistTitle:     "superlist",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{ID: "wishlist1", Title: "superlist", CustomerId: "customer1", Items: []string{"product1"}}, nil)
			},
			expectedError: nil,
		},
		{
			name:              "should return error on product get error",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{ID: "wishlist1", CustomerId: "customer1"}, nil)
				mockProductGetter.EXPECT().Execute(gomock.Any(), "product1").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "should return error on secondt product get error",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1", "product2"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{ID: "wishlist1", CustomerId: "customer1"}, nil)
				mockProductGetter.EXPECT().Execute(gomock.Any(), "product1").Return(&domain.Product{ID: "product1"}, nil)
				mockProductGetter.EXPECT().Execute(gomock.Any(), "product2").Return(nil, errors.New("database error"))

			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "should return not found nil product",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{ID: "wishlist1", CustomerId: "customer1"}, nil)
				mockProductGetter.EXPECT().Execute(gomock.Any(), "product1").Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("product_product1"),
		},
		{
			name:              "should return error on wishlist update error",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{ID: "wishlist1", CustomerId: "customer1"}, nil)
				mockProductGetter.EXPECT().Execute(gomock.Any(), "product1").Return(&domain.Product{ID: "product1"}, nil)
				mockWishlistUpdater.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "should succeed on wishlist update success",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			wishlistTitle:     "superlist",
			products:          []string{"product1"},
			setupMocks: func() {
				mockCustomerGetter.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mockWishlistGetter.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{ID: "wishlist1", CustomerId: "customer1"}, nil)
				mockProductGetter.EXPECT().Execute(gomock.Any(), "product1").Return(&domain.Product{ID: "product1"}, nil)
				mockWishlistUpdater.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			uc := usecase.NewUpdateWishListUseCase(
				mockCustomerGetter,
				mockWishlistGetter,
				mockWishlistUpdater,
				mockProductGetter,
			)

			err := uc.UpdateWishlist(context.Background(), tt.currentCustomerID, &domain.Wishlist{
				ID:         tt.wishlistID,
				Title:      tt.wishlistTitle,
				CustomerId: tt.customerID,
				Items:      tt.products,
			})
			assert.Equal(t, err, tt.expectedError)
		})
	}
}
