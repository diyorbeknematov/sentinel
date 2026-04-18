package apperrors

import "net/http"

func NotFound(msg string, err error) *AppError {
	return New(ErrNotFound, msg, http.StatusNotFound, err)
}

func BadRequest(msg string) *AppError {
	return New(ErrInvalidInput, msg, http.StatusBadRequest, nil)
}

func Internal(err error) *AppError {
	return New(ErrInternal, "internal server error", http.StatusInternalServerError, err)
}

func Conflict(msg string) *AppError {
	return New(ErrAlreadyExists, msg, http.StatusConflict, nil)
}
