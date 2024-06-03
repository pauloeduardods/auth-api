package admin

import (
	"auth-api/src/pkg/app_error"
	"database/sql/driver"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type AdminID uuid.UUID

func ParseAdminID(id string) (AdminID, error) {
	parsedID, err := uuid.Parse(id)
	if err == nil {
		return AdminID(parsedID), nil
	}
	return AdminID(uuid.Nil), app_error.NewApiError(http.StatusBadRequest, "Invalid admin ID", fmt.Sprintf("Field: %s", "ID"))
}

func (id AdminID) String() string {
	return uuid.UUID(id).String()
}

type Admin struct {
	ID    AdminID `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
}

func (id *AdminID) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("scanning a nil value into AdminID")
	}

	switch v := value.(type) {
	case string:
		uid, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = AdminID(uid)
	case []byte:
		uid, err := uuid.Parse(string(v))
		if err != nil {
			return err
		}
		*id = AdminID(uid)
	default:
		return fmt.Errorf("scanning a value of type %T into AdminID", v)
	}

	return nil
}

func (id AdminID) Value() (driver.Value, error) {
	return id.String(), nil
}
