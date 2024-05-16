package service

import (
	"context"
	"monitoring-system/server/pkg/app_error"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type AuthServiceImpl struct {
	client   *cognito.Client
	clientId string
	ctx      context.Context
}

type AuthService interface {
	Login(l LoginInput) (*cognito.InitiateAuthOutput, error)
	SignUp(s SignUpInput) (*cognito.SignUpOutput, error)
	UserInformation(accessToken string) (*cognito.GetUserOutput, error)
	ConfirmSignUp(s ConfirmSignUpInput) (*cognito.ConfirmSignUpOutput, error)
	GetUser(g GetUserInput) (*cognito.GetUserOutput, error)
}

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

type SignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
	Name     string `json:"name" binding:"required" validate:"min=3,max=50"`
}

type ConfirmSignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Code     string `json:"code" binding:"required" validate:"numeric"`
}

type GetUserInput struct {
	AccessToken string `json:"accessToken" binding:"required"`
}

func NewAuthService(ctx context.Context, c *cognito.Client, clientId string) *AuthServiceImpl {
	return &AuthServiceImpl{
		client:   c,
		clientId: clientId,
		ctx:      ctx,
	}
}

func (c *AuthServiceImpl) Login(l LoginInput) (*cognito.InitiateAuthOutput, error) {
	input := &cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": l.Username,
			"PASSWORD": l.Password,
		},
		ClientId: aws.String(c.clientId),
	}
	out, err := c.client.InitiateAuth(c.ctx, input)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return nil, app_error.NewApiError(401, "Invalid username or password")
		}
		if strings.Contains(errorType, "PasswordResetRequiredException") {
			return nil, app_error.NewApiError(401, "Password reset required")
		}
		if strings.Contains(errorType, "UserNotConfirmedException") {
			return nil, app_error.NewApiError(401, "User not confirmed")
		}
		return nil, err
	}
	return out, nil
}

func (c *AuthServiceImpl) SignUp(s SignUpInput) (*cognito.SignUpOutput, error) {
	input := &cognito.SignUpInput{
		ClientId: aws.String(c.clientId),
		Username: aws.String(s.Username),
		Password: aws.String(s.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(s.Name),
			},
		},
	}
	out, err := c.client.SignUp(c.ctx, input)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UsernameExistsException") {
			return nil, app_error.NewApiError(409, "Username already exists")
		}
		return nil, err
	}
	return out, nil
}

func (c *AuthServiceImpl) UserInformation(accessToken string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	return c.client.GetUser(c.ctx, input)
}

func (c *AuthServiceImpl) ConfirmSignUp(s ConfirmSignUpInput) (*cognito.ConfirmSignUpOutput, error) {
	input := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientId),
		Username:         aws.String(s.Username),
		ConfirmationCode: aws.String(s.Code),
	}
	return c.client.ConfirmSignUp(context.Background(), input)
}

func (c *AuthServiceImpl) GetUser(g GetUserInput) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: &g.AccessToken,
	}
	return c.client.GetUser(context.Background(), input)
}
