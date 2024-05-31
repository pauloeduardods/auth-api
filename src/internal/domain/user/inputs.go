package user

import (
	"fmt"
	"monitoring-system/server/src/pkg/app_error"
	"monitoring-system/server/src/pkg/validator"
	"net/http"
	"strings"
)

type CreateUserInput struct {
	ID    UserID
	Name  string
	Email string
	Phone *string
}

func NewCreateUserInput(id, name, email string, phone *string) (CreateUserInput, error) {
	userID, err := ParseUserID(id)
	if err != nil {
		return CreateUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}

	lowerCaseEmail := strings.ToLower(email)
	if err := validator.ValidateEmail(lowerCaseEmail); err != nil {
		return CreateUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Email"))
	}
	if err := validator.ValidateStringLength(name, 3, 100); err != nil {
		return CreateUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
	}
	return CreateUserInput{
		ID:    userID,
		Name:  name,
		Email: lowerCaseEmail,
		Phone: phone,
	}, nil
}

type UpdateUserInput struct {
	ID    UserID
	Name  *string
	Email *string
	Phone *string
}

func NewUpdateUserInput(id string, name, email, phone *string) (UpdateUserInput, error) {
	userID, err := ParseUserID(id)
	if err != nil {
		return UpdateUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}
	if email != nil {
		lowerCaseEmail := strings.ToLower(*email)
		if err := validator.ValidateEmail(lowerCaseEmail); err != nil {
			return UpdateUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Email"))
		}
		email = &lowerCaseEmail
	}
	if name != nil {
		if err := validator.ValidateStringLength(*name, 3, 100); err != nil {
			return UpdateUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
		}
	}
	if phone != nil {
		if err := validator.ValidateStringLength(*phone, 10, 15); err != nil {
			return UpdateUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid phone number length", fmt.Sprintf("Field: %s", "Phone"))
		}
	}

	return UpdateUserInput{
		ID:    userID,
		Name:  name,
		Email: email,
		Phone: phone,
	}, nil
}

type GetUserInput struct {
	ID string
}

func NewGetUserInput(id string) (GetUserInput, error) {
	_, err := ParseUserID(id)
	if err != nil {
		return GetUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}
	return GetUserInput{
		ID: id,
	}, nil
}

type DeleteUserInput struct {
	ID UserID
}

func NewDeleteUserInput(id string) (DeleteUserInput, error) {
	userID, err := ParseUserID(id)
	if err != nil {
		return DeleteUserInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}
	return DeleteUserInput{
		ID: userID,
	}, nil
}

type GetUserByEmailInput struct {
	Email string
}

func NewGetUserByEmailInput(email string) (GetUserByEmailInput, error) {
	lowerCaseEmail := strings.ToLower(email)
	if err := validator.ValidateEmail(lowerCaseEmail); err != nil {
		return GetUserByEmailInput{}, app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Email"))
	}
	return GetUserByEmailInput{
		Email: lowerCaseEmail,
	}, nil
}
