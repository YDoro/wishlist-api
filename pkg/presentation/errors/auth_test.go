package errors_test

import (
	"fmt"
	"testing"

	e "github.com/ydoro/wishlist/pkg/wishlist/presentation/errors"
)

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "isAuthenticationError",
			err:  &e.AuthenticationError{},
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
