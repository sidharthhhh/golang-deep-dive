package errors

import "fmt"

// AppError represents an application error
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Predefined errors
var (
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized access",
		StatusCode: 401,
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "Access forbidden",
		StatusCode: 403,
	}

	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		StatusCode: 404,
	}

	ErrBadRequest = &AppError{
		Code:       "BAD_REQUEST",
		Message:    "Invalid request",
		StatusCode: 400,
	}

	ErrConflict = &AppError{
		Code:       "CONFLICT",
		Message:    "Resource conflict",
		StatusCode: 409,
	}

	ErrInternalServer = &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
		StatusCode: 500,
	}

	ErrRateLimitExceeded = &AppError{
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    "Rate limit exceeded",
		StatusCode: 429,
	}

	ErrTokenExpired = &AppError{
		Code:       "TOKEN_EXPIRED",
		Message:    "Token has expired",
		StatusCode: 401,
	}

	ErrTokenInvalid = &AppError{
		Code:       "TOKEN_INVALID",
		Message:    "Invalid token",
		StatusCode: 401,
	}

	ErrTokenRevoked = &AppError{
		Code:       "TOKEN_REVOKED",
		Message:    "Token has been revoked",
		StatusCode: 401,
	}

	ErrInvalidPassword = &AppError{
		Code:       "INVALID_PASSWORD",
		Message:    "Invalid password",
		StatusCode: 401,
	}
)
