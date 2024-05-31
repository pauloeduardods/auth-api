package user

import (
	"auth-api/src/pkg/app_error"
	"auth-api/src/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

type CreateUserInput struct {
	ID    UserID
	Name  string
	Email string
	Phone *string
}

func (input *CreateUserInput) Validate() error {
	userID, err := ParseUserID(input.ID.String())
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}
	input.ID = userID

	lowerCaseEmail, err := validateEmail(input.Email)
	if err != nil {
		return err
	}
	input.Email = lowerCaseEmail

	if err := validator.ValidateStringLength(input.Name, 3, 100); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
	}

	if input.Phone != nil {
		if err := validator.ValidateStringLength(*input.Phone, 10, 15); err != nil {
			return app_error.NewApiError(http.StatusBadRequest, "Invalid phone number length", fmt.Sprintf("Field: %s", "Phone"))
		}
	}
	return nil
}

type UpdateUserInput struct {
	ID    UserID
	Name  *string
	Email *string
	Phone *string
}

func (input *UpdateUserInput) Validate() error {
	userID, err := ParseUserID(input.ID.String())
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}
	input.ID = userID

	if input.Email != nil {
		lowerCaseEmail, err := validateEmail(*input.Email)
		if err != nil {
			return err
		}
		input.Email = &lowerCaseEmail
	}

	if input.Name != nil {
		if err := validator.ValidateStringLength(*input.Name, 3, 100); err != nil {
			return app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
		}
	}

	if input.Phone != nil {
		if err := validator.ValidateStringLength(*input.Phone, 10, 15); err != nil {
			return app_error.NewApiError(http.StatusBadRequest, "Invalid phone number length", fmt.Sprintf("Field: %s", "Phone"))
		}
	}
	return nil
}

type GetUserInput struct {
	ID string
}

func (input *GetUserInput) Validate() error {
	_, err := ParseUserID(input.ID)
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}
	return nil
}

type DeleteUserInput struct {
	ID UserID
}

func (input *DeleteUserInput) Validate() error {
	userID, err := ParseUserID(input.ID.String())
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
	}
	input.ID = userID
	return nil
}

type GetUserByEmailInput struct {
	Email string
}

func (input *GetUserByEmailInput) Validate() error {
	lowerCaseEmail, err := validateEmail(input.Email)
	if err != nil {
		return err
	}
	input.Email = lowerCaseEmail
	return nil
}

func validateEmail(email string) (string, error) {
	lowerCaseEmail := strings.ToLower(email)
	if err := validator.ValidateEmail(lowerCaseEmail); err != nil {
		return "", app_error.NewApiError(http.StatusBadRequest, "Invalid email format", fmt.Sprintf("Field: %s", "Email"))
	}
	return lowerCaseEmail, nil
}
