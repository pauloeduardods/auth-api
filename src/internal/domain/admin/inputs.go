package admin

import (
	"auth-api/src/pkg/app_error"
	"auth-api/src/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

type CreateAdminInput struct {
	ID     AdminID
	Name   string
	Email  string
	Status AdminStatus
}

func (input *CreateAdminInput) Validate() error {
	adminID, err := ParseAdminID(input.ID.String())
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid admin ID", fmt.Sprintf("Field: %s", "ID"))
	}
	input.ID = adminID

	lowerCaseEmail, err := validateEmail(input.Email)
	if err != nil {
		return err
	}
	input.Email = lowerCaseEmail

	if err := validator.ValidateStringLength(input.Name, 3, 100); err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid name length", fmt.Sprintf("Field: %s", "Name"))
	}

	if input.Status != AdminStatusActive && input.Status != AdminStatusInactive && input.Status != AdminStatusDeleted { //TODO: handle this better
		return app_error.NewApiError(http.StatusBadRequest, "Invalid status", fmt.Sprintf("Field: %s", "Status"))
	}

	return nil
}

type UpdateAdminInput struct {
	ID     AdminID
	Name   *string
	Email  *string
	Status *AdminStatus
}

func (input *UpdateAdminInput) Validate() error {
	adminID, err := ParseAdminID(input.ID.String())
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid admin ID", fmt.Sprintf("Field: %s", "ID"))
	}
	input.ID = adminID

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

	if input.Status != nil && *input.Status != AdminStatusActive && *input.Status != AdminStatusInactive && *input.Status != AdminStatusDeleted { //TODO: handle this better
		return app_error.NewApiError(http.StatusBadRequest, "Invalid status", fmt.Sprintf("Field: %s", "Status"))
	}

	return nil
}

type GetAdminInput struct {
	ID string
}

func (input *GetAdminInput) Validate() error {
	_, err := ParseAdminID(input.ID)
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid admin ID", fmt.Sprintf("Field: %s", "ID"))
	}
	return nil
}

type DeleteAdminInput struct {
	ID AdminID
}

func (input *DeleteAdminInput) Validate() error {
	adminID, err := ParseAdminID(input.ID.String())
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid admin ID", fmt.Sprintf("Field: %s", "ID"))
	}
	input.ID = adminID
	return nil
}

type GetAdminByEmailInput struct {
	Email string
}

func (input *GetAdminByEmailInput) Validate() error {
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

type ChangeStatusAdminInput struct {
	ID     AdminID
	Status AdminStatus
}

func (input *ChangeStatusAdminInput) Validate() error {
	adminID, err := ParseAdminID(input.ID.String())
	if err != nil {
		return app_error.NewApiError(http.StatusBadRequest, "Invalid admin ID", fmt.Sprintf("Field: %s", "ID"))
	}
	input.ID = adminID

	if input.Status != AdminStatusActive && input.Status != AdminStatusInactive && input.Status != AdminStatusDeleted { //TODO: handle this better
		return app_error.NewApiError(http.StatusBadRequest, "Invalid status", fmt.Sprintf("Field: %s", "Status"))
	}

	return nil
}
