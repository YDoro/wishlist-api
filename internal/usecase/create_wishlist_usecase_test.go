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

func TestCreateWishlistUseCase_CreateWishlist(t *testing.T) {
	tests := []struct {
		name              string
		currentCustomerID string
		customerID        string
		title             string
		setupMocks        func(*mocks.MockGetCustomerByIDRepository, *mocks.MockWishlistByTitleRepository, *mocks.MockWishlistCreationRepository, *mocks.MockIDGenerator)
		expectedID        string
		expectedError     error
	}{
		{
			name:              "successful wishlist creation",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetByTitle(gomock.Any(), "customer_123", "Birthday Wishlist").
					Return(nil, nil)

				idGen.EXPECT().
					Generate().
					Return("wishlist_123", nil)

				expectedWishlist := &domain.Wishlist{
					ID:         "wishlist_123",
					CustomerId: "customer_123",
					Title:      "Birthday Wishlist",
					Items:      []string{},
				}

				wishlistCreator.EXPECT().
					Create(gomock.Any(), matchesWishlist(expectedWishlist)).
					Return(nil)
			},
			expectedID:    "wishlist_123",
			expectedError: nil,
		},
		{
			name:              "unauthorized creation attempt",
			currentCustomerID: "different_customer",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
			},
			expectedID:    "",
			expectedError: e.NewUnauthorizedError(),
		},
		{
			name:              "customer not found",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, nil)
			},
			expectedID:    "",
			expectedError: e.NewNotFoundError("customer"),
		},
		{
			name:              "customer getter error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, errors.New("database error"))
			},
			expectedID:    "",
			expectedError: errors.New("database error"),
		},
		{
			name:              "title already in use",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetByTitle(gomock.Any(), "customer_123", "Birthday Wishlist").
					Return(&domain.Wishlist{}, nil)
			},
			expectedID:    "",
			expectedError: &e.ValidationError{Field: "title", Err: "already in use"},
		},
		{
			name:              "wishlist getter error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetByTitle(gomock.Any(), "customer_123", "Birthday Wishlist").
					Return(nil, errors.New("database error"))
			},
			expectedID:    "",
			expectedError: errors.New("database error"),
		},
		{
			name:              "id generation error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetByTitle(gomock.Any(), "customer_123", "Birthday Wishlist").
					Return(nil, nil)

				idGen.EXPECT().
					Generate().
					Return("", errors.New("id generation error"))
			},
			expectedID:    "",
			expectedError: errors.New("id generation error"),
		},
		{
			name:              "wishlist creation error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			title:             "Birthday Wishlist",
			setupMocks: func(customerGetter *mocks.MockGetCustomerByIDRepository,
				wishlistGetter *mocks.MockWishlistByTitleRepository,
				wishlistCreator *mocks.MockWishlistCreationRepository,
				idGen *mocks.MockIDGenerator) {
				customerGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(&domain.Customer{ID: "customer_123"}, nil)

				wishlistGetter.EXPECT().
					GetByTitle(gomock.Any(), "customer_123", "Birthday Wishlist").
					Return(nil, nil)

				idGen.EXPECT().
					Generate().
					Return("wishlist_123", nil)

				wishlistCreator.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedID:    "",
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
			wishlistGetter := mocks.NewMockWishlistByTitleRepository(ctrl)
			wishlistCreator := mocks.NewMockWishlistCreationRepository(ctrl)
			idGen := mocks.NewMockIDGenerator(ctrl)

			tt.setupMocks(customerGetter, wishlistGetter, wishlistCreator, idGen)

			uc := usecase.NewCreateWishlistUseCase(
				wishlistGetter,
				wishlistCreator,
				customerGetter,
				idGen,
			)

			id, err := uc.CreateWishlist(context.Background(), tt.currentCustomerID, tt.customerID, tt.title)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)
		})
	}
}

func matchesWishlist(expected *domain.Wishlist) gomock.Matcher {
	return &wishlistMatcher{expected: expected}
}

type wishlistMatcher struct {
	expected *domain.Wishlist
}

func (m *wishlistMatcher) Matches(x interface{}) bool {
	wishlist, ok := x.(*domain.Wishlist)
	if !ok {
		return false
	}

	return wishlist.ID == m.expected.ID &&
		wishlist.CustomerId == m.expected.CustomerId &&
		wishlist.Title == m.expected.Title &&
		len(wishlist.Items) == len(m.expected.Items)
}

func (m *wishlistMatcher) String() string {
	return "matches wishlist"
}
