package auth

import (
	"auth-api/src/pkg/app_error"
	"auth-api/src/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

type LoginInput struct {
	Username string
	Password string
}

func (input *LoginInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername

	if err := validator.ValidatePassword(input.Password); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid password", fmt.Sprintf("Field: %s", "Password"))
	}
	return nil
}

type SignUpInput struct {
	Username string
	Password string
	Name     string
}

func (input *SignUpInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername

	if err := validator.ValidatePassword(input.Password); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid password", fmt.Sprintf("Field: %s", "Password"))
	}

	if err := validator.ValidateStringLength(input.Name, 3, 50); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
	}
	return nil
}

type ConfirmSignUpInput struct {
	Username string
	Code     string
}

func (input *ConfirmSignUpInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername

	if err := validator.ValidateNumeric(input.Code); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid code", fmt.Sprintf("Field: %s", "Code"))
	}
	return nil
}

type GetMeInput struct {
	AccessToken string
}

func (input *GetMeInput) Validate() error {
	if len(input.AccessToken) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Access token is required", fmt.Sprintf("Field: %s", "AccessToken"))
	}
	return nil
}

type RefreshTokenInput struct {
	RefreshToken string
}

func (input *RefreshTokenInput) Validate() error {
	if len(input.RefreshToken) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Refresh token is required", fmt.Sprintf("Field: %s", "RefreshToken"))
	}
	return nil
}

type CreateAdminInput struct {
	Password string
	Name     string
	Username string
}

func (input *CreateAdminInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername

	if err := validator.ValidatePassword(input.Password); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid password", fmt.Sprintf("Field: %s", "Password"))
	}

	if err := validator.ValidateStringLength(input.Name, 3, 50); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
	}
	return nil
}

type AddGroupInput struct {
	Username  string
	GroupName UserGroup
}

func (input *AddGroupInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername

	if input.GroupName != Admin && input.GroupName != User {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid user group", fmt.Sprintf("Field: %s", "GroupName"))
	}
	return nil
}

type RemoveGroupInput struct {
	Username  string
	GroupName UserGroup
}

func (input *RemoveGroupInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername

	if input.GroupName != Admin && input.GroupName != User {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid user group", fmt.Sprintf("Field: %s", "GroupName"))
	}
	return nil
}

type AddMFAInput struct {
	AccessToken string
}

func (input *AddMFAInput) Validate() error {
	if len(input.AccessToken) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Access token is required", fmt.Sprintf("Field: %s", "AccessToken"))
	}
	return nil
}

type VerifyMFAInput struct {
	Code     string
	Username string
	Session  string
}

func (input *VerifyMFAInput) Validate() error {
	if err := validator.ValidateNumeric(input.Code); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid code", fmt.Sprintf("Field: %s", "Code"))
	}

	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername

	if len(input.Session) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Session is required", fmt.Sprintf("Field: %s", "Session"))
	}
	return nil
}

type AdminRemoveMFAInput struct {
	Username string
}

func (input *AdminRemoveMFAInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername
	return nil
}

type RemoveMFAInput struct {
	AccessToken string
}

func (input *RemoveMFAInput) Validate() error {
	if len(input.AccessToken) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Access token is required", fmt.Sprintf("Field: %s", "AccessToken"))
	}
	return nil
}

type DeleteUserInput struct {
	Username string
}

func (input *DeleteUserInput) Validate() error {
	lowerCaseUsername, err := validateEmail(input.Username)
	if err != nil {
		return err
	}
	input.Username = lowerCaseUsername
	return nil
}

type ActivateMFAInput struct {
	AccessToken string
	Code        string
}

func (input *ActivateMFAInput) Validate() error {
	if len(input.AccessToken) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Access token is required", fmt.Sprintf("Field: %s", "AccessToken"))
	}

	if err := validator.ValidateNumeric(input.Code); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid code", fmt.Sprintf("Field: %s", "Code"))
	}
	return nil
}

type LogoutInput struct {
	AccessToken string
}

func (input *LogoutInput) Validate() error {
	if len(input.AccessToken) == 0 {
		return app_error.NewApiError(http.StatusBadRequest, "Access token is required", fmt.Sprintf("Field: %s", "AccessToken"))
	}
	return nil
}

func validateEmail(username string) (string, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return "", app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	return lowerCaseUsername, nil
}
