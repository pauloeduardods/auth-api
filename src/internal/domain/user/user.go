package user

import (
	"database/sql/driver"
	"fmt"
	"monitoring-system/server/src/pkg/app_error"
	"net/http"

	"github.com/google/uuid"
)

type UserID uuid.UUID

func ParseUserID(id string) (UserID, error) {
	parsedID, err := uuid.Parse(id)
	if err == nil {
		return UserID(parsedID), nil
	}
	return UserID(uuid.Nil), app_error.NewApiError(http.StatusBadRequest, "Invalid user ID", fmt.Sprintf("Field: %s", "ID"))
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
