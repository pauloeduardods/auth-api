package auth

import (
	"fmt"
	"monitoring-system/server/src/pkg/app_error"
	"monitoring-system/server/src/pkg/validator"
	"net/http"
	"strings"
)

type LoginInput struct {
	Username string
	Password string
}

func NewLoginInput(username, password string) (LoginInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return LoginInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	if err := validator.ValidatePassword(password); err != nil {
		return LoginInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid password", fmt.Sprintf("Field: %s", "Password"))
	}
	return LoginInput{
		Username: lowerCaseUsername,
		Password: password,
	}, nil
}

type SignUpInput struct {
	Username string
	Password string
	Name     string
}

func NewSignUpInput(username, password, name string) (SignUpInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return SignUpInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	if err := validator.ValidatePassword(password); err != nil {
		return SignUpInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid password", fmt.Sprintf("Field: %s", "Password"))
	}
	if err := validator.ValidateStringLength(name, 3, 50); err != nil {
		return SignUpInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
	}
	return SignUpInput{
		Username: lowerCaseUsername,
		Password: password,
		Name:     name,
	}, nil
}

type ConfirmSignUpInput struct {
	Username string
	Code     string
}

func NewConfirmSignUpInput(username, code string) (ConfirmSignUpInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return ConfirmSignUpInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	if err := validator.ValidateNumeric(code); err != nil {
		return ConfirmSignUpInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid code", fmt.Sprintf("Field: %s", "Code"))
	}
	return ConfirmSignUpInput{
		Username: lowerCaseUsername,
		Code:     code,
	}, nil
}

type GetMeInput struct {
	AccessToken string
}

func NewGetMeInput(accessToken string) (GetMeInput, error) {
	if len(accessToken) == 0 {
		return GetMeInput{}, app_error.NewApiError(http.StatusBadRequest, "Access token is required", fmt.Sprintf("Field: %s", "AccessToken"))
	}
	return GetMeInput{
		AccessToken: accessToken,
	}, nil
}

type RefreshTokenInput struct {
	RefreshToken string
}

func NewRefreshTokenInput(refreshToken string) (RefreshTokenInput, error) {
	if len(refreshToken) == 0 {
		return RefreshTokenInput{}, app_error.NewApiError(http.StatusBadRequest, "Refresh token is required", fmt.Sprintf("Field: %s", "RefreshToken"))
	}
	return RefreshTokenInput{
		RefreshToken: refreshToken,
	}, nil
}

type CreateAdminInput struct {
	Password string
	Name     string
	Username string
}

func NewCreateAdminInput(username, password, name string) (CreateAdminInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return CreateAdminInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	if err := validator.ValidatePassword(password); err != nil {
		return CreateAdminInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid password", fmt.Sprintf("Field: %s", "Password"))
	}
	if err := validator.ValidateStringLength(name, 3, 50); err != nil {
		return CreateAdminInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
	}
	return CreateAdminInput{
		Password: password,
		Name:     name,
		Username: lowerCaseUsername,
	}, nil
}

type AddGroupInput struct {
	Username  string
	GroupName UserGroup
}

func NewAddGroupInput(username string, groupName UserGroup) (AddGroupInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return AddGroupInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	if groupName != Admin && groupName != User {
		return AddGroupInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid user group", fmt.Sprintf("Field: %s", "GroupName"))
	}
	return AddGroupInput{
		Username:  lowerCaseUsername,
		GroupName: groupName,
	}, nil
}

type RemoveGroupInput struct {
	Username  string
	GroupName UserGroup
}

func NewRemoveGroupInput(username string, groupName UserGroup) (RemoveGroupInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return RemoveGroupInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	if groupName != Admin && groupName != User {
		return RemoveGroupInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid user group", fmt.Sprintf("Field: %s", "GroupName"))
	}
	return RemoveGroupInput{
		Username:  lowerCaseUsername,
		GroupName: groupName,
	}, nil
}

type AddMFAInput struct {
	AccessToken string
}

func NewAddMFAInput(accessToken string) (AddMFAInput, error) {
	if len(accessToken) == 0 {
		return AddMFAInput{}, app_error.NewApiError(http.StatusBadRequest, "Access token is required", fmt.Sprintf("Field: %s", "AccessToken"))
	}
	return AddMFAInput{
		AccessToken: accessToken,
	}, nil
}

type VerifyMFAInput struct {
	Code     string
	Username string
	Session  string
}

func NewVerifyMFAInput(code, username, session string) (VerifyMFAInput, error) {
	if err := validator.ValidateNumeric(code); err != nil {
		return VerifyMFAInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid code", fmt.Sprintf("Field: %s", "Code"))
	}
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return VerifyMFAInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	if len(session) == 0 {
		return VerifyMFAInput{}, app_error.NewApiError(http.StatusBadRequest, "Session is required", fmt.Sprintf("Field: %s", "Session"))
	}
	return VerifyMFAInput{
		Code:     code,
		Username: lowerCaseUsername,
		Session:  session,
	}, nil
}

type RemoveMFAInput struct {
	Username string
}

func NewRemoveMFAInput(username string) (RemoveMFAInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return RemoveMFAInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	return RemoveMFAInput{
		Username: lowerCaseUsername,
	}, nil
}

type DeleteUserInput struct {
	Username string
}

func NewDeleteUserInput(username string) (DeleteUserInput, error) {
	lowerCaseUsername := strings.ToLower(username)
	if err := validator.ValidateEmail(lowerCaseUsername); err != nil {
		return DeleteUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Username"))
	}
	return DeleteUserInput{
		Username: lowerCaseUsername,
	}, nil
}
