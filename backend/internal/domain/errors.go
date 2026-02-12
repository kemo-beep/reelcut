package domain

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrValidation          = errors.New("validation error")
	ErrConflict            = errors.New("conflict")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
	ErrInternal            = errors.New("internal error")
	ErrInsufficientCredits = errors.New("insufficient credits")
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
