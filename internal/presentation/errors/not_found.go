package errors

import "fmt"

type NotFoundError struct {
	Resource string
}

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
	}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
