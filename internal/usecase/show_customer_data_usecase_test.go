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

func TestGetCustomerData_ShowCustomerData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGetCustomerByIDRepository(ctrl)

	tests := []struct {
		name              string
		currentCustomerID string
		customerID        string
		setupMocks        func()
		expectedCustomer  *domain.OutgoingCustomer
		expectedError     error
	}{
		{
			name:              "unauthorized access - different customer ID",
			currentCustomerID: "current_id",
			customerID:        "different_id",
			setupMocks:        func() {},
			expectedCustomer:  nil,
			expectedError:     &e.UnauthorizedError{},
		},
		{
			name:              "successful customer data retrieval",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			setupMocks: func() {
				customer := &domain.Customer{
					ID:        "customer_123",
					Name:      "John Doe",
					Email:     "john@example.com",
					Password:  "hashed_password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				mockRepo.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)
			},
			expectedCustomer: &domain.OutgoingCustomer{
				ID:    "customer_123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expectedError: nil,
		},
		{
			name:              "customer not found",
			currentCustomerID: "id_123",
			customerID:        "id_123",
			setupMocks: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "id_123").
					Return(nil, nil)
			},
			expectedCustomer: nil,
			expectedError:    &e.NotFoundError{Resource: "customer"},
		},
		{
			name:              "repository error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			setupMocks: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, errors.New("database error"))
			},
			expectedCustomer: nil,
			expectedError:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			uc := usecase.NewGetCustomerData(mockRepo)
			customer, err := uc.ShowCustomerData(context.Background(), tt.currentCustomerID, tt.customerID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, customer)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCustomer.ID, customer.ID)
				assert.Equal(t, tt.expectedCustomer.Name, customer.Name)
				assert.Equal(t, tt.expectedCustomer.Email, customer.Email)
				assert.NotZero(t, customer.CreatedAt)
			}
		})
	}
}
