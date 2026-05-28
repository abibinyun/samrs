package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response standar untuk sukses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Response standar untuk sukses dengan metadata
type APIResponseWithMeta struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Response standar untuk error
type ErrResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse memberikan output JSON sukses yang seragam
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessResponseWithMeta memberikan output JSON sukses dengan metadata
func SuccessResponseWithMeta(c *gin.Context, message string, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, APIResponseWithMeta{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse memberikan output JSON error yang seragam
func ErrorResponse(c *gin.Context, code int, message string, errs interface{}) {
	c.JSON(code, ErrResponse{
		Success: false,
		Message: message,
		Errors:  normalizeErrors(errs),
	})
}

func normalizeErrors(errs interface{}) interface{} {
	if errs == nil {
		return nil
	}

	switch v := errs.(type) {
	case error:
		return map[string]interface{}{"detail": v.Error()}
	case string:
		if v == "" {
			return nil
		}
		return map[string]interface{}{"detail": v}
	case []string:
		return map[string]interface{}{"details": v}
	case map[string]interface{}:
		return v
	case gin.H:
		return v
	default:
		return errs
	}
}
