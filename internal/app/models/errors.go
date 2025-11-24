package models

import "errors"

// Common errors used across models
var (
	ErrInsufficientStock       = errors.New("insufficient stock available")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrNotFound                = errors.New("resource not found")
	ErrUnauthorized            = errors.New("unauthorized access")
	ErrInvalidInput            = errors.New("invalid input")
	ErrDuplicateEntry          = errors.New("duplicate entry")
	ErrNegativeStock           = errors.New("stock cannot be negative")
	ErrInvalidQuantity         = errors.New("quantity must be greater than zero")
)
