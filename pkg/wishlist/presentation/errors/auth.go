package errors

import (
	"fmt"

	"github.com/ydoro/wishlist/internal/customer/domain"
)

type AuthenticationError struct {
	AuthType domain.AuthMethod `json:"auth_type"`
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("AuthenticationError: failed to authenticate using %s authentication method", e.AuthType)
}

func IsAuthenticationError(err error) bool {
	if _, ok := err.(*AuthenticationError); ok {
		return true
	}
	return false
}

func NewAuthenticationError(authtype domain.AuthMethod) error {
	return &AuthenticationError{
		AuthType: authtype,
	}
}
