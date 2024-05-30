package auth

import "context"

type UserGroup string

const (
	Admin UserGroup = "Admin"
	User  UserGroup = "User"
)

type AuthService interface {
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
	SignUp(ctx context.Context, input SignUpInput) (*SignUpOutput, error)
	RollbackSignUp(ctx context.Context, input SignUpInput) error
	ConfirmSignUp(ctx context.Context, input ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetMe(ctx context.Context, input GetMeInput) (*GetMeOutput, error)
	ValidateToken(ctx context.Context, token string) (*Claims, error)
	AddGroup(ctx context.Context, input AddGroupInput) error
	RemoveGroup(ctx context.Context, input RemoveGroupInput) error
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error)
	CreateAdmin(ctx context.Context, input CreateAdminInput) (*CreateAdminOutput, error)
	AddMFA(ctx context.Context, input AddMFAInput) (*AddMFAOutput, error)
	ActivateMFA(ctx context.Context, input ActivateMFAInput) error
	VerifyMFA(ctx context.Context, input VerifyMFAInput) (*LoginOutput, error)
	RemoveMFA(ctx context.Context, input RemoveMFAInput) error
}

type AuthClient interface {
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
	SignUp(ctx context.Context, input SignUpInput) (*SignUpOutput, error)
	DeleteUser(ctx context.Context, input DeleteUserInput) error
	ConfirmSignUp(ctx context.Context, input ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetMe(ctx context.Context, input GetMeInput) (*GetMeOutput, error)
	ValidateToken(ctx context.Context, token string) (*Claims, error)
	AddGroup(ctx context.Context, input AddGroupInput) error
	RemoveGroup(ctx context.Context, input RemoveGroupInput) error
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error)
	CreateAdmin(ctx context.Context, input CreateAdminInput) (*CreateAdminOutput, error)
	AddMFA(ctx context.Context, input AddMFAInput) (*AddMFAOutput, error)
	ActivateMFA(ctx context.Context, input ActivateMFAInput) error
	VerifyMFA(ctx context.Context, input VerifyMFAInput) (*LoginOutput, error)
	RemoveMFA(ctx context.Context, input RemoveMFAInput) error
}

type Claims struct {
	Email      string   `json:"email"`
	Id         string   `json:"id"`
	UserGroups []string `json:"groups"`
}
