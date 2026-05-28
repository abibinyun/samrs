package util

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// HTTPStatusFromError maps common error types to HTTP status codes.
// Default is 500 to avoid masking unexpected/internal errors.
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case IsErrorKind(err, ErrKindValidation):
		return http.StatusBadRequest
	case IsErrorKind(err, ErrKindNotFound):
		return http.StatusNotFound
	case IsErrorKind(err, ErrKindConflict):
		return http.StatusConflict
	case IsErrorKind(err, ErrKindForbidden):
		return http.StatusForbidden
	case IsErrorKind(err, ErrKindUnauthorized):
		return http.StatusUnauthorized
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return http.StatusBadRequest
	}

	switch err.(type) {
	case *json.SyntaxError, *json.UnmarshalTypeError:
		return http.StatusBadRequest
	case *strconv.NumError, *time.ParseError:
		return http.StatusBadRequest
	}
	if strings.Contains(strings.ToLower(err.Error()), "invalid uuid") {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

// ErrorResponseFromErr sends standardized error response with mapped status.
func ErrorResponseFromErr(c *gin.Context, message string, err error) {
	ErrorResponse(c, HTTPStatusFromError(err), message, errorDetailsForResponse(err))
}

func errorDetailsForResponse(err error) interface{} {
	if err == nil {
		return nil
	}
	if IsErrorKind(err, ErrKindValidation) ||
		IsErrorKind(err, ErrKindNotFound) ||
		IsErrorKind(err, ErrKindConflict) ||
		IsErrorKind(err, ErrKindForbidden) ||
		IsErrorKind(err, ErrKindUnauthorized) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		return err
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return err
	}

	switch err.(type) {
	case *json.SyntaxError, *json.UnmarshalTypeError:
		return err
	case *strconv.NumError, *time.ParseError:
		return err
	}
	if strings.Contains(strings.ToLower(err.Error()), "invalid uuid") {
		return err
	}

	return nil
}
