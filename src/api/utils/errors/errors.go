package errors

import "net/http"

type ApiError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	EStatus  int    `json:"status"`
	EMessage string `json:"message"`
	EError   string `json:"error,omitempty"`
}

func (e *apiError) Status() int {
	return e.EStatus
}
func (e *apiError) Message() string {
	return e.EMessage
}
func (e *apiError) Error() string {
	return e.EError
}

func NewNotFoundApiError(message string) ApiError {
	return &apiError{
		EStatus:  http.StatusNotFound,
		EMessage: message,
	}
}

func NewInternalServerError(message string) ApiError {
	return &apiError{
		EStatus:  http.StatusInternalServerError,
		EMessage: message,
	}
}

func NewBadRequestError(message string) ApiError {
	return &apiError{
		EStatus:  http.StatusBadRequest,
		EMessage: message,
	}
}

func NewApiError(statusCode int, message string) ApiError {
	return &apiError{
		EStatus:  statusCode,
		EMessage: message,
	}
}
