package auth

import (
	"auth-api/src/pkg/app_error"
	"fmt"
)

var (
	ErrInvalidGroup               = app_error.NewApiError(400, "Invalid group", fmt.Sprintf("Field: %s", "Group"))
	ErrInvalidMfaCode             = app_error.NewApiError(400, "Invalid MFA code")
	ErrInvalidAccessCode          = app_error.NewApiError(401, "Invalid access token")
	ErrFailedToVerifySoftwareMfa  = app_error.NewApiError(400, "Failed to verify software MFA")
	ErrFailedToRespondToChallenge = app_error.NewApiError(400, "Failed to respond to challenge")
	ErrInvalidUsernameOrPassword  = app_error.NewApiError(401, "Invalid username or password")
	ErrPasswordResetRequired      = app_error.NewApiError(401, "Password reset required")
	ErrUserNotConfirmed           = app_error.NewApiError(401, "User not confirmed")
	ErrUserAlreadyExists          = app_error.NewApiError(409, "User already exists")
	ErrInvalidRefreshToken        = app_error.NewApiError(401, "Invalid refresh token")
	ErrUserNotFound               = app_error.NewApiError(404, "User not found")
	ErrUserAlreadyConfirmed       = app_error.NewApiError(409, "User already confirmed")
	ErrInvalidUserStatus          = app_error.NewApiError(400, "Invalid user status")
)
