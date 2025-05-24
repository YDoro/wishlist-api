package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/domain"
	"github.com/ydoro/wishlist/internal/usecase"
	mocks "github.com/ydoro/wishlist/mock/domain"
	e "github.com/ydoro/wishlist/pkg/presentation/errors"
	"github.com/ydoro/wishlist/pkg/presentation/inputs"
	"go.uber.org/mock/gomock"
)

func TestPasswordAuthenticationUseCase_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHashComparer := mocks.NewMockHashComparer(ctrl)
	mockUserGetter := mocks.NewMockGetCustomerByEmailRepository(ctrl)
	mockEncrypter := mocks.NewMockEncrypter(ctrl)

	testUser := &domain.Customer{
		ID:       "user123",
		Email:    "test@example.com",
		Password: "hashedPassword",
		Name:     "Test User",
	}

	testUserJSON, _ := json.Marshal(testUser)

	tests := []struct {
		name          string
		credentials   inputs.PwdAuth
		setupMocks    func()
		expectedToken string
		expectedError error
	}{
		{
			name: "successful authentication",
			credentials: inputs.PwdAuth{
				Email:    "test@example.com",
				Password: "correctPassword",
			},
			setupMocks: func() {
				mockUserGetter.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(testUser, nil)

				mockHashComparer.EXPECT().
					Compare("hashedPassword", "correctPassword").
					Return(nil)

				mockEncrypter.EXPECT().
					Encrypt(string(testUserJSON)).
					Return("valid.jwt.token", nil)
			},
			expectedToken: "valid.jwt.token",
			expectedError: nil,
		},
		{
			name: "user not found",
			credentials: inputs.PwdAuth{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			setupMocks: func() {
				mockUserGetter.EXPECT().
					GetByEmail(gomock.Any(), "nonexistent@example.com").
					Return(nil, nil)
			},
			expectedToken: "",
			expectedError: e.NewAuthenticationError(domain.AuthMethodPassword),
		},
		{
			name: "database error",
			credentials: inputs.PwdAuth{
				Email:    "test@example.com",
				Password: "password",
			},
			setupMocks: func() {
				mockUserGetter.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(nil, errors.New("database error"))
			},
			expectedToken: "",
			expectedError: errors.New("database error"),
		},
		{
			name: "incorrect password",
			credentials: inputs.PwdAuth{
				Email:    "test@example.com",
				Password: "wrongPassword",
			},
			setupMocks: func() {
				mockUserGetter.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(testUser, nil)

				mockHashComparer.EXPECT().
					Compare("hashedPassword", "wrongPassword").
					Return(errors.New("hash comparison failed"))
			},
			expectedToken: "",
			expectedError: e.NewAuthenticationError(domain.AuthMethodPassword),
		},
		{
			name: "token generation error",
			credentials: inputs.PwdAuth{
				Email:    "test@example.com",
				Password: "correctPassword",
			},
			setupMocks: func() {
				mockUserGetter.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(testUser, nil)

				mockHashComparer.EXPECT().
					Compare("hashedPassword", "correctPassword").
					Return(nil)

				mockEncrypter.EXPECT().
					Encrypt(string(testUserJSON)).
					Return("", errors.New("encryption failed"))
			},
			expectedToken: "",
			expectedError: errors.New("encryption failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			useCase := usecase.NewPasswordAuthenticationUseCase(
				mockHashComparer,
				mockUserGetter,
				mockEncrypter,
			)

			token, err := useCase.Authenticate(context.Background(), tt.credentials)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}
