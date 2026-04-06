package errutil

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	StatusCode    int     `json:"statusCode"`
	ExceptionCode string  `json:"exceptionCode"`
	Message       string  `json:"message"`
	Errors        []Error `json:"errors,omitempty"`
}

// Error represents a structured error detail.
type Error struct {
	Code    string   `json:"code"`
	Message []string `json:"message"`
}

// RecoverFromPanic recovers from a panic and logs the error.
func RecoverFromPanic() {
	if recover := recover(); recover != nil {
		log.Println("recovered from ", recover)
	}
}

// NewAPIError creates a standardized error response.
func NewAPIError(statusCode int, errCode, message string, errorDetails []Error) *ErrorResponse {
	return &ErrorResponse{
		StatusCode:    statusCode,
		ExceptionCode: errCode,
		Message:       message,
		Errors:        errorDetails,
	}
}

// Panic throws an ErrorResponse as a panic to be caught by ErrorMiddleware.
func Panic(statusCode int, errCode, message string, errorDetails []Error) {
	panic(NewAPIError(statusCode, errCode, message, errorDetails))
}

// Common pre-built errors

func BadRequest(message string) *ErrorResponse {
	return NewAPIError(http.StatusBadRequest, "BAD_REQUEST", message, []Error{})
}

func Unauthorized(message string) *ErrorResponse {
	return NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", message, []Error{})
}

func NotFound(message string) *ErrorResponse {
	return NewAPIError(http.StatusNotFound, "NOT_FOUND", message, []Error{})
}

func NotImplemented(message string) *ErrorResponse {
	return NewAPIError(http.StatusNotImplemented, "NOT_IMPLEMENTED", message, []Error{})
}

func Internal(message string) *ErrorResponse {
	return NewAPIError(http.StatusInternalServerError, "INTERNAL_ERROR", message, []Error{})
}

// ErrorMiddleware handles panics and converts them to JSON error responses.
func ErrorMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				var errResp *ErrorResponse

				switch e := rec.(type) {
				case *ErrorResponse:
					errResp = e
				default:
					errResp = &ErrorResponse{
						StatusCode:    http.StatusInternalServerError,
						ExceptionCode: "INTERNAL_ERROR",
						Message:       "Internal error occurred",
						Errors:        []Error{},
					}
				}

				context.JSON(errResp.StatusCode, errResp)
				context.Abort()
			}
		}()

		context.Next()
	}
}
