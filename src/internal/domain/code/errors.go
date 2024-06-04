package code

import "auth-api/src/pkg/app_error"

var (
	ErrInvalidCode  = app_error.NewApiError(400, "invalid_code", "Invalid code")
	ErrCodeExpired  = app_error.NewApiError(400, "code_expired", "Code expired")
	ErrCodeNotFound = app_error.NewApiError(404, "code_not_found", "Code not found")
)
