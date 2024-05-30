package auth

import "context"

type auth struct {
	cognitoAuth AuthClient
}

func NewAuthService(cognitoAuth AuthClient) AuthService {
	return &auth{
		cognitoAuth: cognitoAuth,
	}
}

func (a *auth) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	return a.cognitoAuth.Login(ctx, input)
}

func (a *auth) SignUp(ctx context.Context, input SignUpInput) (*SignUpOutput, error) {
	return a.cognitoAuth.SignUp(ctx, input)
}

func (a *auth) RollbackSignUp(ctx context.Context, input SignUpInput) error {
	return a.cognitoAuth.DeleteUser(ctx, DeleteUserInput{Username: input.Username})
}

func (a *auth) ConfirmSignUp(ctx context.Context, input ConfirmSignUpInput) (*ConfirmSignUpOutput, error) {
	return a.cognitoAuth.ConfirmSignUp(ctx, input)
}

func (a *auth) GetMe(ctx context.Context, input GetMeInput) (*GetMeOutput, error) {
	return a.cognitoAuth.GetMe(ctx, input)
}

func (a *auth) ValidateToken(ctx context.Context, token string) (*Claims, error) {
	return a.cognitoAuth.ValidateToken(ctx, token)
}

func (a *auth) AddGroup(ctx context.Context, input AddGroupInput) error {
	return a.cognitoAuth.AddGroup(ctx, input)
}

func (a *auth) RemoveGroup(ctx context.Context, input RemoveGroupInput) error {
	return a.cognitoAuth.RemoveGroup(ctx, input)
}

func (a *auth) RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	return a.cognitoAuth.RefreshToken(ctx, input)
}

func (a *auth) CreateAdmin(ctx context.Context, input CreateAdminInput) (*CreateAdminOutput, error) {
	return a.cognitoAuth.CreateAdmin(ctx, input)
}

func (a *auth) AddMFA(ctx context.Context, input AddMFAInput) (*AddMFAOutput, error) {
	return a.cognitoAuth.AddMFA(ctx, input)
}

func (a *auth) ActivateMFA(ctx context.Context, input ActivateMFAInput) error {
	return a.cognitoAuth.ActivateMFA(ctx, input)
}

func (a *auth) VerifyMFA(ctx context.Context, input VerifyMFAInput) (*LoginOutput, error) {
	return a.cognitoAuth.VerifyMFA(ctx, input)
}

func (a *auth) RemoveMFA(ctx context.Context, input RemoveMFAInput) error {
	return a.cognitoAuth.RemoveMFA(ctx, input)
}
