package utils

import (
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Common error types
var (
	ErrValidationFailed = &AppError{
		Code:    "VALIDATION_FAILED",
		Message: "Validation failed",
	}
	ErrNotFound = &AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
	}
	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized access",
	}
	ErrForbidden = &AppError{
		Code:    "FORBIDDEN",
		Message: "Access forbidden",
	}
	ErrConflict = &AppError{
		Code:    "CONFLICT",
		Message: "Resource conflict",
	}
	ErrInternalServer = &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
	}
)

// HTTPStatus returns the HTTP status code for the error
func HTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Code {
		case "VALIDATION_FAILED":
			return http.StatusBadRequest
		case "NOT_FOUND":
			return http.StatusNotFound
		case "UNAUTHORIZED":
			return http.StatusUnauthorized
		case "FORBIDDEN":
			return http.StatusForbidden
		case "CONFLICT":
			return http.StatusConflict
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *AppError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	Page     int   `json:"page,omitempty"`
	PageSize  int   `json:"page_size,omitempty"`
	Total    int64 `json:"total,omitempty"`
	TotalPages int  `json:"total_pages,omitempty"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(err error) *Response {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		appErr = ErrInternalServer
		appErr.Details = err.Error()
	}
	return &Response{
		Success: false,
		Error:   appErr,
	}
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, meta *Meta) *Response {
	return &Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}
