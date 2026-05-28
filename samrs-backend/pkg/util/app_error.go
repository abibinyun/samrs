package util

import "errors"

const (
	ErrKindValidation  = "validation"
	ErrKindNotFound    = "not_found"
	ErrKindConflict    = "conflict"
	ErrKindForbidden   = "forbidden"
	ErrKindUnauthorized = "unauthorized"
)

type AppError struct {
	Kind    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func ErrValidation(message string) error {
	return AppError{Kind: ErrKindValidation, Message: message}
}

func ErrNotFound(message string) error {
	return AppError{Kind: ErrKindNotFound, Message: message}
}

func ErrConflict(message string) error {
	return AppError{Kind: ErrKindConflict, Message: message}
}

func ErrForbidden(message string) error {
	return AppError{Kind: ErrKindForbidden, Message: message}
}

func ErrUnauthorized(message string) error {
	return AppError{Kind: ErrKindUnauthorized, Message: message}
}

func IsErrorKind(err error, kind string) bool {
	var appErr AppError
	if errors.As(err, &appErr) {
		return appErr.Kind == kind
	}
	return false
}
