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

func TestDeleteCustomerUseCase_DeleteCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockUpdateCustomerRepository(ctrl)
	mockGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)

	tests := []struct {
		name              string
		currentCustomerID string
		customerID        string
		setupMocks        func()
		expectedError     error
	}{
		{
			name:              "successful deletion",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			setupMocks: func() {
				customer := &domain.Customer{
					ID:        "customer_123",
					Name:      "Test User",
					Email:     "test@example.com",
					Password:  "hashed_password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)

				mockUpdater.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, c *domain.Customer) error {
						assert.NotZero(t, c.DeletedAt)
						return nil
					})
			},
			expectedError: nil,
		},
		{
			name:              "unauthorized deletion attempt",
			currentCustomerID: "different_id",
			customerID:        "customer_123",
			setupMocks:        func() {},
			expectedError:     e.NewUnauthorizedError(),
		},
		{
			name:              "customer not found",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			setupMocks: func() {
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("customer"),
		},
		{
			name:              "database error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			setupMocks: func() {
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:              "update error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			setupMocks: func() {
				customer := &domain.Customer{
					ID:        "customer_123",
					Name:      "Test User",
					Email:     "test@example.com",
					Password:  "hashed_password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)

				mockUpdater.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			uc := usecase.NewDeleteCustomerUseCase(mockUpdater, mockGetter)
			err := uc.DeleteCustomer(context.Background(), tt.currentCustomerID, tt.customerID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
