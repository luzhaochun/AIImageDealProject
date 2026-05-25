package apperr

import "errors"

var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrNoImageAvailable   = errors.New("no image available")
	ErrInvalidTransition  = errors.New("invalid status transition")
	ErrValidation         = errors.New("validation error")
	ErrDiscardLimit       = errors.New("discard limit exceeded")
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
