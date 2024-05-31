package auth

type UserGroup string

const (
	Admin UserGroup = "Admin"
	User  UserGroup = "User"
)

type Claims struct {
	Email      string   `json:"email"`
	Id         string   `json:"id"`
	UserGroups []string `json:"groups"`
}
