package admin

import "auth-api/src/pkg/app_error"

var (
	ErrAdminNotFound      = app_error.NewApiError(404, "Admin not found", "Field: id")
	ErrAdminAlreadyExists = app_error.NewApiError(409, "Admin already exists", "Field: email")
)
