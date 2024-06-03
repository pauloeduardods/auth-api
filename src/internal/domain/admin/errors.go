package admin

import "auth-api/src/pkg/app_error"

var (
	ErrAdminNotFound         = app_error.NewApiError(404, "Admin not found", "Field: id")
	ErrAdminAlreadyExists    = app_error.NewApiError(409, "Admin already exists", "Field: email")
	ErrAdminStatusNotChanged = app_error.NewApiError(400, "Admin status not changed", "Field: status")
	ErrAdminDeleted          = app_error.NewApiError(400, "Admin deleted", "Field: id")
)
