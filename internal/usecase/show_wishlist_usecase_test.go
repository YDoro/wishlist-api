package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/usecase"
	mocks "github.com/ydoro/wishlist/mock/domain"
	"go.uber.org/mock/gomock"
)

func TestShowWishlist_AuthorizationAndValidation(t *testing.T) {
	tests := []struct {
		name              string
		currentCustomerID string
		customerID        string
		wishlistID        string
		setupMocks        func(*mocks.MockGetCustomerByIDRepository, *mocks.MockWishlistByIdRepository)
		expectedError     error
	}{
		{
			name:              "should return unauthorized when current customer is different",
			currentCustomerID: "customer1",
			customerID:        "customer2",
			wishlistID:        "wishlist1",
			setupMocks:        func(mc *mocks.MockGetCustomerByIDRepository, mw *mocks.MockWishlistByIdRepository) {},
			expectedError:     &e.UnauthorizedError{},
		},
		{
			name:              "should return not found when customer does not exist",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			setupMocks: func(mc *mocks.MockGetCustomerByIDRepository, mw *mocks.MockWishlistByIdRepository) {
				mc.EXPECT().GetByID(gomock.Any(), "customer1").Return(nil, nil)
			},
			expectedError: &e.NotFoundError{Resource: "customer"},
		},
		{
			name:              "should return not found when wishlist does not exist",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			setupMocks: func(mc *mocks.MockGetCustomerByIDRepository, mw *mocks.MockWishlistByIdRepository) {
				mc.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mw.EXPECT().GetById(gomock.Any(), "wishlist1").Return(nil, nil)
			},
			expectedError: &e.NotFoundError{Resource: "wishlist"},
		},
		{
			name:              "should return unauthorized when wishlist belongs to another customer",
			currentCustomerID: "customer1",
			customerID:        "customer1",
			wishlistID:        "wishlist1",
			setupMocks: func(mc *mocks.MockGetCustomerByIDRepository, mw *mocks.MockWishlistByIdRepository) {
				mc.EXPECT().GetByID(gomock.Any(), "customer1").Return(&domain.Customer{ID: "customer1"}, nil)
				mw.EXPECT().GetById(gomock.Any(), "wishlist1").Return(&domain.Wishlist{
					ID:         "wishlist1",
					CustomerId: "customer2",
				}, nil)
			},
			expectedError: &e.UnauthorizedError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
			mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
			mockProductGetter := mocks.NewMockGetProductUseCase(ctrl)

			tt.setupMocks(mockCustomerGetter, mockWishlistGetter)

			uc := usecase.NewShowWishlistUseCase(mockWishlistGetter, mockCustomerGetter, mockProductGetter)
			wishlist, err := uc.ShowWishlist(context.Background(), tt.currentCustomerID, tt.customerID, tt.wishlistID)

			assert.Error(t, err)
			assert.IsType(t, tt.expectedError, err)
			assert.Nil(t, wishlist)
		})
	}
}

func TestShowWishlist_EmptyWishlist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := &domain.Customer{
		ID:        "customer1",
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	emptyWishlist := &domain.Wishlist{
		ID:         "wishlist1",
		CustomerId: "customer1",
		Title:      "Empty Wishlist",
		Items:      []string{},
	}

	mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
	mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
	mockProductGetter := mocks.NewMockGetProductUseCase(ctrl)

	mockCustomerGetter.EXPECT().
		GetByID(gomock.Any(), "customer1").
		Return(customer, nil)

	mockWishlistGetter.EXPECT().
		GetById(gomock.Any(), "wishlist1").
		Return(emptyWishlist, nil)

	uc := usecase.NewShowWishlistUseCase(mockWishlistGetter, mockCustomerGetter, mockProductGetter)
	result, err := uc.ShowWishlist(context.Background(), "customer1", "customer1", "wishlist1")

	assert.NoError(t, err)
	assert.Equal(t, emptyWishlist.ID, result.ID)
	assert.Equal(t, emptyWishlist.Title, result.Title)
	assert.Equal(t, customer, result.Customer)
	assert.Empty(t, result.Items)
}

func TestShowWishlist_WithProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := &domain.Customer{
		ID:        "customer1",
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	wishlist := &domain.Wishlist{
		ID:         "wishlist1",
		CustomerId: "customer1",
		Title:      "My Wishlist",
		Items:      []string{"product1", "product2"},
	}

	product1 := &domain.Product{ID: "product1", Name: "Product 1"}
	product2 := &domain.Product{ID: "product2", Name: "Product 2"}

	mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
	mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
	mockProductGetter := mocks.NewMockGetProductUseCase(ctrl)

	mockCustomerGetter.EXPECT().
		GetByID(gomock.Any(), "customer1").
		Return(customer, nil)

	mockWishlistGetter.EXPECT().
		GetById(gomock.Any(), "wishlist1").
		Return(wishlist, nil)

	// Product fetch expectations with any order
	mockProductGetter.EXPECT().
		Execute(gomock.Any(), "product1").
		Return(product1, nil)

	mockProductGetter.EXPECT().
		Execute(gomock.Any(), "product2").
		Return(product2, nil)

	uc := usecase.NewShowWishlistUseCase(mockWishlistGetter, mockCustomerGetter, mockProductGetter)
	result, err := uc.ShowWishlist(context.Background(), "customer1", "customer1", "wishlist1")

	assert.NoError(t, err)
	assert.Equal(t, wishlist.ID, result.ID)
	assert.Equal(t, wishlist.Title, result.Title)
	assert.Equal(t, customer, result.Customer)
	assert.Len(t, result.Items, 2)

	foundProducts := make(map[string]bool)
	for _, item := range result.Items {
		foundProducts[item.ID] = true
		if item.ID == "product1" {
			assert.Equal(t, "Product 1", item.Name)
		} else if item.ID == "product2" {
			assert.Equal(t, "Product 2", item.Name)
		}
	}
	assert.True(t, foundProducts["product1"])
	assert.True(t, foundProducts["product2"])
}

func TestShowWishlist_RepositoryErrors(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*mocks.MockGetCustomerByIDRepository, *mocks.MockWishlistByIdRepository)
		expectedError error
	}{
		{
			name: "should return error when customer repository fails",
			setupMocks: func(mc *mocks.MockGetCustomerByIDRepository, mw *mocks.MockWishlistByIdRepository) {
				mc.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "should return error when wishlist repository fails",
			setupMocks: func(mc *mocks.MockGetCustomerByIDRepository, mw *mocks.MockWishlistByIdRepository) {
				mc.EXPECT().
					GetByID(gomock.Any(), "customer1").
					Return(&domain.Customer{ID: "customer1"}, nil)
				mw.EXPECT().
					GetById(gomock.Any(), "wishlist1").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
			mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
			mockProductGetter := mocks.NewMockGetProductUseCase(ctrl)

			tt.setupMocks(mockCustomerGetter, mockWishlistGetter)

			uc := usecase.NewShowWishlistUseCase(mockWishlistGetter, mockCustomerGetter, mockProductGetter)
			result, err := uc.ShowWishlist(context.Background(), "customer1", "customer1", "wishlist1")

			assert.Error(t, err)
			assert.Equal(t, tt.expectedError.Error(), err.Error())
			assert.Nil(t, result)
		})
	}
}

func TestShowWishlist_ProductErrors(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*mocks.MockGetProductUseCase)
		expectedError error
	}{
		{
			name: "should not return error when product service fails",
			setupMocks: func(mp *mocks.MockGetProductUseCase) {
				mp.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("product service error"))
			},
			expectedError: nil,
		},
		{
			name: "should handle context cancellation",
			setupMocks: func(mp *mocks.MockGetProductUseCase) {
				mp.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, context.Canceled)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customer := &domain.Customer{ID: "customer1"}
			wishlist := &domain.Wishlist{
				ID:         "wishlist1",
				CustomerId: "customer1",
				Items:      []string{"product1"},
			}

			mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
			mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
			mockProductGetter := mocks.NewMockGetProductUseCase(ctrl)

			mockCustomerGetter.EXPECT().
				GetByID(gomock.Any(), "customer1").
				Return(customer, nil)

			mockWishlistGetter.EXPECT().
				GetById(gomock.Any(), "wishlist1").
				Return(wishlist, nil)

			tt.setupMocks(mockProductGetter)

			uc := usecase.NewShowWishlistUseCase(mockWishlistGetter, mockCustomerGetter, mockProductGetter)
			result, err := uc.ShowWishlist(context.Background(), "customer1", "customer1", "wishlist1")

			if tt.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			}

		})
	}
}

func TestShowWishlist_ContextCancellationInGoroutine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customer := &domain.Customer{ID: "customer1"}
	wishlist := &domain.Wishlist{
		ID:         "wishlist1",
		CustomerId: "customer1",
		Items:      []string{"product1", "product2"},
	}

	mockCustomerGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
	mockWishlistGetter := mocks.NewMockWishlistByIdRepository(ctrl)
	mockProductGetter := mocks.NewMockGetProductUseCase(ctrl)

	mockCustomerGetter.EXPECT().
		GetByID(gomock.Any(), "customer1").
		Return(customer, nil)

	mockWishlistGetter.EXPECT().
		GetById(gomock.Any(), "wishlist1").
		Return(wishlist, nil)

	mockProductGetter.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string) (*domain.Product, error) {
			cancel() // Cancel context before product fetch
			<-ctx.Done()
			return nil, ctx.Err()
		}).AnyTimes()

	uc := usecase.NewShowWishlistUseCase(mockWishlistGetter, mockCustomerGetter, mockProductGetter)
	result, err := uc.ShowWishlist(ctx, "customer1", "customer1", "wishlist1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Items)
}
