package common

import (
	"fmt"
)

type BaseError struct {
	Message string
	Cause   error
}

func (e BaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e BaseError) Unwrap() error {
	return e.Cause
}

type AuthorizationError struct{ BaseError }

func NewAuthorizationError(message string, cause error) AuthorizationError {
	return AuthorizationError{BaseError{Message: message, Cause: cause}}
}

type ValidationError struct{ BaseError }

func NewValidationError(message string, cause error) ValidationError {
	return ValidationError{BaseError{Message: message, Cause: cause}}
}

type UnprocessableEntityError struct{ BaseError }

func NewUnprocessableEntityError(message string, cause error) UnprocessableEntityError {
	return UnprocessableEntityError{BaseError{Message: message, Cause: cause}}
}
