package apperrors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrAlreadyExists  ErrorCode = "ALREADY_EXISTS"
	ErrInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrInternal       ErrorCode = "INTERNAL"
	ErrNoRowsAffected ErrorCode = "no rows affected"
)

type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == string(code)
	}
	return false
}

func New(code ErrorCode, msg string, status int, err error) *AppError {
	return &AppError{
		Code:    string(code),
		Message: msg,
		Status:  status,
		Err:     err,
	}
}
