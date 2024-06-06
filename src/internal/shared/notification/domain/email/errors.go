package email

import "auth-api/src/pkg/app_error"

var (
	ErrEmailToEmpty      = app_error.NewApiError(400, "Email to is empty")
	ErrEmailSubjectEmpty = app_error.NewApiError(400, "Email subject is empty")
	ErrEmailBodyEmpty    = app_error.NewApiError(400, "Email body is empty")
)
