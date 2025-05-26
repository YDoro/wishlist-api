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

func TestCreateCustomerUseCase_CreateCustomerWithEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerCreationRepository(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)

	tests := []struct {
		name          string
		input         domain.IncommingCustomer
		setupMocks    func()
		expectedID    string
		expectedError error
	}{
		{
			name: "successful customer creation",
			input: domain.IncommingCustomer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockHasher.EXPECT().
					Hash("password123").
					Return("hashed_password", nil)

				mockIDGen.EXPECT().
					Generate().
					Return("generated_id", nil)

				expectedCustomer := &domain.Customer{
					ID:        "generated_id",
					Name:      "John Doe",
					Email:     "john@example.com",
					Password:  "hashed_password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				mockRepo.EXPECT().
					Create(gomock.Any(), matchesCustomer(expectedCustomer)).
					Return(nil)
			},
			expectedID:    "generated_id",
			expectedError: nil,
		},
		{
			name: "empty password",
			input: domain.IncommingCustomer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "",
			},
			setupMocks:    func() {},
			expectedID:    "",
			expectedError: e.NewRequiredFieldError("password"),
		},
		{
			name: "password hashing fails",
			input: domain.IncommingCustomer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockHasher.EXPECT().
					Hash("password123").
					Return("", errors.New("hashing error"))
			},
			expectedID:    "",
			expectedError: errors.New("hashing error\nfailed to hash password"),
		},
		{
			name: "id generation fails",
			input: domain.IncommingCustomer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockHasher.EXPECT().
					Hash("password123").
					Return("hashed_password", nil)

				mockIDGen.EXPECT().
					Generate().
					Return("", errors.New("id generation error"))
			},
			expectedID:    "",
			expectedError: errors.New("id generation error\nfailed to generate customer ID"),
		},
		{
			name: "repository creation fails",
			input: domain.IncommingCustomer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockHasher.EXPECT().
					Hash("password123").
					Return("hashed_password", nil)

				mockIDGen.EXPECT().
					Generate().
					Return("generated_id", nil)

				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("repository error"))
			},
			expectedID:    "generated_id",
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			uc := usecase.NewCreateCustomerUseCase(mockRepo, mockIDGen, mockHasher)
			id, err := uc.CreateCustomerWithEmail(context.Background(), tt.input)

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

// matchesCustomer is a custom matcher for comparing Customer objects
func matchesCustomer(expected *domain.Customer) gomock.Matcher {
	return &customerMatcher{expected: expected}
}

type customerMatcher struct {
	expected *domain.Customer
}

func (m *customerMatcher) Matches(x interface{}) bool {
	customer, ok := x.(*domain.Customer)
	if !ok {
		return false
	}

	// Compare all fields except timestamps
	return customer.ID == m.expected.ID &&
		customer.Name == m.expected.Name &&
		customer.Email == m.expected.Email &&
		customer.Password == m.expected.Password
}

func (m *customerMatcher) String() string {
	return "matches customer"
}
