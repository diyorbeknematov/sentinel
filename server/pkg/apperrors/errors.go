package apperrors

import (
	"errors"
	"fmt"
	"log/slog"
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
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
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

func LogByError(log *slog.Logger, msg string, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		// AppError emas — oddiy xato, ERROR level
		log.Error(msg, "err", err)
		return
	}

	switch appErr.Code {
	case string(ErrInternal):
		log.Error(msg, "err", appErr.Err, "code", appErr.Code)
	case string(ErrNotFound):
		log.Warn(msg, "err", appErr.Err, "code", appErr.Code)
	case string(ErrInvalidInput):
		log.Warn(msg, "err", appErr.Err, "code", appErr.Code)
	case string(ErrAlreadyExists):
		log.Info(msg, "err", appErr.Err, "code", appErr.Code)
	default:
		log.Error(msg, "err", appErr.Err, "code", appErr.Code)
	}
}
