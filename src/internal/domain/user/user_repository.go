package user

import "monitoring-system/server/src/pkg/app_error"

var (
	ErrUserNotFound      = app_error.NewApiError(404, "User not found", "Field: id")
	ErrUserAlreadyExists = app_error.NewApiError(409, "User already exists", "Field: email")
)

type UserRepository interface {
	GetByID(input *GetUserInput) (*User, error)
	GetByEmail(email *GetUserByEmailInput) (*User, error)
	Create(input *CreateUserInput) error
	Update(user *UpdateUserInput) error
	Delete(id *DeleteUserInput) error
}
