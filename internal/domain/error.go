package domain

import "errors"

// ErrorResponse defines a standard error response.
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

var (
	ErrScoreNotFound = errors.New("score not found")
)
