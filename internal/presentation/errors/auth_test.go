package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	e "github.com/ydoro/wishlist/internal/presentation/errors"
)

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "isAuthenticationError",
			err:  e.NewAuthenticationError("password"),
			want: true,
		},
		{
			name: "isNotAuthenticationError",
			err:  fmt.Errorf("some other error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.IsAuthenticationError(tt.err); got != tt.want {
				t.Errorf("IsAuthenticationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "AuthenticationError message",
			err:  e.NewAuthenticationError("password"),
			want: "password authentication method",
		},
		{
			name: "not AuthenticationError",
			err:  fmt.Errorf("some other error"),
			want: "some other error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, tt.err.Error(), tt.want)
		})
	}
}
