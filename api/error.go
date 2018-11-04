package api

import (
	"net/http"
)

// ErrorResponse is the API error response format.
type ErrorResponse struct {
	Success bool   `json:"success,omitempty"`
	Code    int    `json:"code,omitempty"`
	Error   string `json:"error,omitempty"`
}

// This holds the default errors.
var (
	ErrInternal     = NewError(http.StatusInternalServerError, "An internal server error has occured.")
	ErrBadRequest   = NewError(http.StatusBadRequest, "Request could not be processed. Bad request.")
	ErrMissingFile  = NewError(http.StatusBadRequest, "The file is missing.")
	ErrFileTooLarge = NewError(http.StatusRequestEntityTooLarge, "The file is too large.")
	ErrFileType     = NewError(http.StatusUnprocessableEntity, "The file type is not supported.")
)

// NewError returns a new ErrorResponse
func NewError(code int, err string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Code:    code,
		Error:   err,
	}
}
