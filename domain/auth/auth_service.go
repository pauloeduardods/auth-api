package auth

type auth struct {
	cognitoAuth CognitoAuth
}

func New(cognitoAuth CognitoAuth) Auth {
	return &auth{
		cognitoAuth: cognitoAuth,
	}
}

func (a *auth) Login(input LoginInput) (*LoginOutput, error) {
	return a.cognitoAuth.Login(input)
}

func (a *auth) SignUp(input SignUpInput) (*SignUpOutput, error) {
	return a.cognitoAuth.SignUp(input)
}

func (a *auth) ConfirmSignUp(input ConfirmSignUpInput) (*ConfirmSignUpOutput, error) {
	return a.cognitoAuth.ConfirmSignUp(input)
}

func (a *auth) GetUser(input GetUserInput) (*GetUserOutput, error) {
	return a.cognitoAuth.GetUser(input)
}

func (a *auth) ValidateToken(token string) (*Claims, error) {
	return a.cognitoAuth.ValidateToken(token)
}
