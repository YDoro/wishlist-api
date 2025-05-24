package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	e "github.com/ydoro/wishlist/pkg/presentation/errors"
)

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "ValidationError",
			err:  e.NewRequiredFieldError("email"),
			want: true,
		},
		{
			name: "not ValidationError",
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
			if got := e.IsValidationError(tt.err); got != tt.want {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ValidationError",
			err:  e.NewRequiredFieldError("email"),
			want: "'email' is required",
		},
		{
			name: "not ValidationError",
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
