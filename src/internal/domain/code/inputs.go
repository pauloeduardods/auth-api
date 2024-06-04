package code

import (
	"auth-api/src/pkg/app_error"
	"fmt"
	"net/http"
	"time"
)

type GenerateAndSaveInput struct {
	Identifier        string
	ExpiresAt         time.Time
	Length            int
	CanContainLetters bool
}

func (input *GenerateAndSaveInput) Validate() error {
	if len(input.Identifier) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Identifier is required", fmt.Sprintf("Field: %s", "Identifier"))
	}
	if input.ExpiresAt.IsZero() {
		return app_error.NewApiError(http.StatusBadRequest, "ExpiresAt is required", fmt.Sprintf("Field: %s", "ExpiresAt"))
	}
	if input.Length <= 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Length is required", fmt.Sprintf("Field: %s", "Length"))
	}

	return nil
}

type VerifyCodeInput struct {
	Identifier string
	Code       string
}

func (input *VerifyCodeInput) Validate() error {
	if len(input.Identifier) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Identifier is required", fmt.Sprintf("Field: %s", "Identifier"))
	}
	if len(input.Code) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Code is required", fmt.Sprintf("Field: %s", "Code"))
	}
	return nil
}
