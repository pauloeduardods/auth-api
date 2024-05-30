package user

import (
	"database/sql/driver"
	"fmt"
	"monitoring-system/server/src/pkg/app_error"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = app_error.NewApiError(404, "User not found", "Field: id")
	ErrUserAlreadyExists = app_error.NewApiError(409, "User already exists", "Field: email")
)

type UserID uuid.UUID

func ParseUserID(id string) (UserID, error) {
	parsedID, err := uuid.Parse(id)
	return UserID(parsedID), err
}

func (id UserID) String() string {
	return uuid.UUID(id).String()
}

type User struct {
	ID    UserID  `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Phone *string `json:"phone,omitempty"`
}

type UserService interface {
	GetByID(input *GetUserInput) (*User, error)
	GetByEmail(email *GetUserByEmailInput) (*User, error)
	Create(input *CreateUserInput) error
	RollbackCreate(input *CreateUserInput) error
	Update(user *UpdateUserInput) (backup *User, err error)
	RollbackUpdate(backup *User) error
	Delete(id *DeleteUserInput) error
}

type UserRepository interface {
	GetByID(input *GetUserInput) (*User, error)
	GetByEmail(email *GetUserByEmailInput) (*User, error)
	Create(input *CreateUserInput) error
	Update(user *UpdateUserInput) error
	Delete(id *DeleteUserInput) error
}

func (id *UserID) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("scanning a nil value into UserID")
	}

	switch v := value.(type) {
	case string:
		uid, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = UserID(uid)
	case []byte:
		uid, err := uuid.Parse(string(v))
		if err != nil {
			return err
		}
		*id = UserID(uid)
	default:
		return fmt.Errorf("scanning a value of type %T into UserID", v)
	}

	return nil
}

func (id UserID) Value() (driver.Value, error) {
	return id.String(), nil
}
