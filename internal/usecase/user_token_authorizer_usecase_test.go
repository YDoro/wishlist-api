package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/usecase"
	mocks "github.com/ydoro/wishlist/mock/domain"
	"go.uber.org/mock/gomock"
)

func TestUserTokenAuthorizerUseCase_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDecrypter := mocks.NewMockDecrypter(ctrl)

	tests := []struct {
		name          string
		token         string
		setupMocks    func()
		expectedID    string
		expectedError error
	}{
		{
			name:  "successful token authorization",
			token: "valid.jwt.token",
			setupMocks: func() {
				mockDecrypter.EXPECT().
					Decrypt("valid.jwt.token").
					Return("user_123", nil)
			},
			expectedID:    "user_123",
			expectedError: nil,
		},
		{
			name:  "invalid token",
			token: "invalid.token",
			setupMocks: func() {
				mockDecrypter.EXPECT().
					Decrypt("invalid.token").
					Return("", errors.New("invalid token"))
			},
			expectedID:    "",
			expectedError: errors.New("invalid token"),
		},
		{
			name:  "empty token",
			token: "",
			setupMocks: func() {
				mockDecrypter.EXPECT().
					Decrypt("").
					Return("", errors.New("empty token"))
			},
			expectedID:    "",
			expectedError: errors.New("empty token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			uc := usecase.NewUserTokenAuthorizerUseCase(mockDecrypter)
			userID, err := uc.Authorize(context.Background(), tt.token)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}
