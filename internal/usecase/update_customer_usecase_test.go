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

func TestUpdateCustomerUseCase_UpdateCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockUpdateCustomerRepository(ctrl)
	mockGetter := mocks.NewMockGetCustomerByIDRepository(ctrl)
	mockEmailGetter := mocks.NewMockGetCustomerByEmailRepository(ctrl)

	getBaseCustomer := func() *domain.Customer {
		return &domain.Customer{
			ID:        "customer_123",
			Name:      "Old Name",
			Email:     "old@example.com",
			Password:  "hashed_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	tests := []struct {
		name              string
		currentCustomerID string
		customerID        string
		updateData        domain.CustomerEditableFields
		setupMocks        func()
		expectedCustomer  *domain.OutgoingCustomer
		expectedError     error
	}{
		{
			name:              "successful update with name and email",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			updateData: domain.CustomerEditableFields{
				Name:  "New Name",
				Email: "new@example.com",
			},
			setupMocks: func() {
				customer := getBaseCustomer()
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)

				mockEmailGetter.EXPECT().
					GetByEmail(gomock.Any(), "new@example.com").
					Return(nil, nil)

				updatedCustomer := *customer
				updatedCustomer.Name = "New Name"
				updatedCustomer.Email = "new@example.com"

				mockUpdater.EXPECT().
					Update(gomock.Any(), matchesCustomer(&updatedCustomer)).
					Return(nil)
			},
			expectedCustomer: &domain.OutgoingCustomer{
				ID:    "customer_123",
				Name:  "New Name",
				Email: "new@example.com",
			},
			expectedError: nil,
		},
		{
			name:              "unauthorized update attempt",
			currentCustomerID: "different_id",
			customerID:        "customer_123",
			updateData: domain.CustomerEditableFields{
				Name: "New Name",
			},
			setupMocks:       func() {},
			expectedCustomer: nil,
			expectedError:    e.NewUnauthorizedError(),
		},
		{
			name:              "customer not found",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			updateData: domain.CustomerEditableFields{
				Name: "New Name",
			},
			setupMocks: func() {
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(nil, errors.New("customer not found"))
			},
			expectedCustomer: nil,
			expectedError:    errors.New("customer not found"),
		},
		{
			name:              "error retrieving customer by email",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			updateData: domain.CustomerEditableFields{
				Email: "existing@example.com",
			},
			setupMocks: func() {
				customer := getBaseCustomer()
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)

				mockEmailGetter.EXPECT().
					GetByEmail(gomock.Any(), "existing@example.com").
					Return(&domain.Customer{ID: "other_customer"}, errors.New("database error"))
			},
			expectedCustomer: nil,
			expectedError:    errors.New("database error"),
		},
		{
			name:              "email already in use",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			updateData: domain.CustomerEditableFields{
				Email: "existing@example.com",
			},
			setupMocks: func() {
				customer := getBaseCustomer()
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)

				mockEmailGetter.EXPECT().
					GetByEmail(gomock.Any(), "existing@example.com").
					Return(&domain.Customer{ID: "other_customer"}, nil)
			},
			expectedCustomer: nil,
			expectedError:    &e.ValidationError{Field: "email", Err: "email already in use"},
		},
		{
			name:              "successful partial update - name only",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			updateData: domain.CustomerEditableFields{
				Name: "New Name",
			},
			setupMocks: func() {
				customer := getBaseCustomer()
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)

				updatedCustomer := *customer
				updatedCustomer.Name = "New Name"

				mockUpdater.EXPECT().
					Update(gomock.Any(), matchesCustomer(&updatedCustomer)).
					Return(nil)
			},
			expectedCustomer: &domain.OutgoingCustomer{
				ID:    "customer_123",
				Name:  "New Name",
				Email: "old@example.com",
			},
			expectedError: nil,
		},
		{
			name:              "repository update error",
			currentCustomerID: "customer_123",
			customerID:        "customer_123",
			updateData: domain.CustomerEditableFields{
				Name: "New Name",
			},
			setupMocks: func() {
				customer := getBaseCustomer()
				mockGetter.EXPECT().
					GetByID(gomock.Any(), "customer_123").
					Return(customer, nil)

				mockUpdater.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedCustomer: nil,
			expectedError:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			uc := usecase.NewUpdateCustomerUseCase(mockUpdater, mockGetter, mockEmailGetter)
			customer, err := uc.UpdateCustomer(context.Background(), tt.currentCustomerID, tt.customerID, tt.updateData)

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
