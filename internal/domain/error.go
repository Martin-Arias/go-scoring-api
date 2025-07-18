package domain

import "errors"

// ErrorResponse defines a standard error response.
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

var (
	ErrScoreNotFound = errors.New("score not found")
	ErrGameNotFound  = errors.New("game not found")
	ErrUserNotFound  = errors.New("user not found")

	ErrGameAlreadyExists     = errors.New("game with the same name already exists")
	ErrUsernameAlreadyExists = errors.New("user with the same username already exists")

	ErrGameCreation   = errors.New("error creating game")
	ErrFetchingUsers  = errors.New("error fetching users")
	ErrCreatingScores = errors.New("error creating initial scores")
)
