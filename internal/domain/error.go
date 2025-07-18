package domain

// ErrorResponse defines a standard error response.
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
