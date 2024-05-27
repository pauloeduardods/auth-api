package auth

type Auth interface {
	Login(LoginInput) (*LoginOutput, error)
	SignUp(SignUpInput) (*SignUpOutput, error)
	// UserInformation(accessToken string) (*cognito.GetUserOutput, error)
	ConfirmSignUp(ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetUser(GetUserInput) (*GetUserOutput, error)
	ValidateToken(token string) (*Claims, error)
}

type CognitoAuth interface {
	Login(LoginInput) (*LoginOutput, error)
	SignUp(SignUpInput) (*SignUpOutput, error)
	// UserInformation(accessToken string) (*cognito.GetUserOutput, error)
	ConfirmSignUp(ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetUser(GetUserInput) (*GetUserOutput, error)
	ValidateToken(token string) (*Claims, error)
}

type LoginOutput struct {
	AccessToken  string `json:"accessToken"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

type SignUpOutput struct {
	IsConfirmed bool `json:"isConfirmed"`
}

type SignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
	Name     string `json:"name" binding:"required" validate:"min=3,max=50"`
}

type ConfirmSignUpOutput struct {
}

type ConfirmSignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Code     string `json:"code" binding:"required" validate:"numeric"`
}

type GetUserOutput struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type GetUserInput struct {
	AccessToken string `json:"accessToken" form:"accessToken" binding:"required"`
}

type Claims struct {
	Email string `json:"email"`
	Id    string `json:"id"`
}
