package errors_test

import (
	"fmt"
	"testing"

	e "github.com/ydoro/wishlist/pkg/wishlist/presentation/errors"
)

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "ValidationError",
			err:  &e.ValidationError{},
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
