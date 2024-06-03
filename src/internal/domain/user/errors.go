package user

import "auth-api/src/pkg/app_error"

var (
	ErrUserNotFound      = app_error.NewApiError(404, "User not found", "Field: id")
	ErrUserAlreadyExists = app_error.NewApiError(409, "User already exists", "Field: email")
)
