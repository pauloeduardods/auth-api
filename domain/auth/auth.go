package auth

type UserGroup string

const (
	Admin UserGroup = "Admin"
	User  UserGroup = "User"
)

type Auth interface {
	Login(LoginInput) (*LoginOutput, error)
	SignUp(SignUpInput) (*SignUpOutput, error)
	ConfirmSignUp(ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetUser(GetUserInput) (*GetUserOutput, error)
	ValidateToken(token string) (*Claims, error)
	AddGroup(AddGroupInput) error
	RemoveGroup(RemoveGroupInput) error
	RefreshToken(RefreshTokenInput) (*RefreshTokenOutput, error)
	CreateAdmin(CreateAdminInput) (*CreateAdminOutput, error)
	AddMFA(AddMFAInput) (*AddMFAOutput, error)
	VerifyMFA(VerifyMFAInput) (*LoginOutput, error)
	RemoveMFA(RemoveMFAInput) error
}

type AuthClient interface {
	Login(LoginInput) (*LoginOutput, error)
	SignUp(SignUpInput) (*SignUpOutput, error)
	ConfirmSignUp(ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetUser(GetUserInput) (*GetUserOutput, error)
	ValidateToken(token string) (*Claims, error)
	AddGroup(AddGroupInput) error
	RemoveGroup(RemoveGroupInput) error
	RefreshToken(RefreshTokenInput) (*RefreshTokenOutput, error)
	CreateAdmin(CreateAdminInput) (*CreateAdminOutput, error)
	AddMFA(AddMFAInput) (*AddMFAOutput, error)
	VerifyMFA(VerifyMFAInput) (*LoginOutput, error)
	RemoveMFA(RemoveMFAInput) error
}

type Claims struct {
	Email      string   `json:"email"`
	Id         string   `json:"id"`
	UserGroups []string `json:"groups"`
}
