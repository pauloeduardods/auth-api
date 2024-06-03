package auth

import (
	"auth-api/src/pkg/app_error"
	"fmt"
)

var (
	ErrInvalidGroup = app_error.NewApiError(400, "Invalid group", fmt.Sprintf("Field: %s", "Group"))
)
