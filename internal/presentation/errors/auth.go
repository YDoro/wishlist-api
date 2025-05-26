package errors

import (
	"fmt"

	"github.com/ydoro/wishlist/internal/domain"
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

type UnauthorizedError struct {
}

func (e *UnauthorizedError) Error() string {
	return "UnauthorizedError: user is not authorized to access this resource"
}
func NewUnauthorizedError() error {
	return &UnauthorizedError{}
}

func IsUnauthorizedError(err error) bool {
	if _, ok := err.(*UnauthorizedError); ok {
		return true
	}
	return false
}
