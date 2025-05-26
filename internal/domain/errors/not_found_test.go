package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/domain/errors"
)

func TestNotFoundError(t *testing.T) {
	err := errors.NewNotFoundError("customer")
	assert.Equal(t, "customer not found", err.Error())
	assert.True(t, errors.IsNotFoundError(err))
	assert.False(t, errors.IsNotFoundError(&errors.ValidationError{Field: "test", Err: "test"}))
}
