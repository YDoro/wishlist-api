package errors

import "fmt"

type ValidationError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

// IsValidationError checks if the given error is a ValidationError.
func IsValidationError(err error) bool {
	if _, ok := err.(*ValidationError); ok {
		return true
	}
	return false
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("ValidationError field '%s' %s", e.Field, e.Err)
}

func NewRequiredFieldError(field string) error {
	return &ValidationError{
		Field: field,
		Err:   "is required",
	}
}
