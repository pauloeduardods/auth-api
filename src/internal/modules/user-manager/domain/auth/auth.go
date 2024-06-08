package auth

type UserGroup string

const (
	GroupAdmin UserGroup = "Admin"
	GroupUser  UserGroup = "User"
)

type UserStatus string

const (
	Unconfirmed       UserStatus = "UNCONFIRMED"
	Confirmed         UserStatus = "CONFIRMED"
	Unknown           UserStatus = "UNKNOWN"
	ResetRequired     UserStatus = "RESET_REQUIRED"
	ForceChangePasswd UserStatus = "FORCE_CHANGE_PASSWORD"
)

type Claims struct {
	Email      string   `json:"email"`
	Id         string   `json:"id"`
	UserGroups []string `json:"groups"`
}

type User struct {
	Id     string     `json:"id"`
	Email  string     `json:"email"`
	Name   string     `json:"name"`
	Status UserStatus `json:"status"`
}

func (us *UserStatus) Scan(value interface{}) error {
	if value == nil {
		return ErrInvalidUserStatus
	}

	switch v := value.(type) {
	case string:
		*us = UserStatus(v)
	case []byte:
		*us = UserStatus(string(v))
	case *string:
		if v != nil {
			*us = UserStatus(*v)
		} else {
			return ErrInvalidUserStatus
		}
	default:
		return ErrInvalidUserStatus
	}

	return nil
}
