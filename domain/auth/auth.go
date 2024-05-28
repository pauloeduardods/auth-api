package auth

import "strings"

type UserGroup string

const (
	Admin UserGroup = "Admin"
	User  UserGroup = "User"
)

type Auth interface {
	Login(LoginInput) (*LoginOutput, error)
	SignUp(SignUpInput) (*SignUpOutput, error)
	// UserInformation(accessToken string) (*cognito.GetUserOutput, error)
	ConfirmSignUp(ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetUser(GetUserInput) (*GetUserOutput, error)
	ValidateToken(token string) (*Claims, error)
	AddGroup(AddGroupInput) error
	RemoveGroup(RemoveGroupInput) error
}

type CognitoAuth interface {
	Login(LoginInput) (*LoginOutput, error)
	SignUp(SignUpInput) (*SignUpOutput, error)
	// UserInformation(accessToken string) (*cognito.GetUserOutput, error)
	ConfirmSignUp(ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetUser(GetUserInput) (*GetUserOutput, error)
	ValidateToken(token string) (*Claims, error)
	AddGroup(AddGroupInput) error
	RemoveGroup(RemoveGroupInput) error
}

type Claims struct {
	Email      string   `json:"email"`
	Id         string   `json:"id"`
	UserGroups []string `json:"groups"`
}

type LoginInput struct {
	Username string
	Password string
}

func NewLoginInput(username, password string) LoginInput {
	lowerCaseUsername := strings.ToLower(username)
	return LoginInput{
		Username: lowerCaseUsername,
		Password: password,
	}
}

type LoginOutput struct {
	AccessToken  string `json:"accessToken"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignUpInput struct {
	Username  string
	Password  string
	Name      string
	GroupName UserGroup
}

func NewSignUpInput(username, password, name string, groupName UserGroup) SignUpInput {
	lowerCaseUsername := strings.ToLower(username)
	return SignUpInput{
		Username:  lowerCaseUsername,
		Password:  password,
		Name:      name,
		GroupName: groupName,
	}
}

type SignUpOutput struct {
	IsConfirmed bool `json:"isConfirmed"`
}

type ConfirmSignUpInput struct {
	Username string
	Code     string
}

func NewConfirmSignUpInput(username, code string) ConfirmSignUpInput {
	lowerCaseUsername := strings.ToLower(username)
	return ConfirmSignUpInput{
		Username: lowerCaseUsername,
		Code:     code,
	}
}

type ConfirmSignUpOutput struct {
}

type GetUserInput struct {
	AccessToken string
}

func NewGetUserInput(accessToken string) GetUserInput {
	return GetUserInput{
		AccessToken: accessToken,
	}
}

type GetUserOutput struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type AddGroupInput struct {
	Username  string
	GroupName UserGroup
}

type RemoveGroupInput struct {
	Username  string
	GroupName UserGroup
}
