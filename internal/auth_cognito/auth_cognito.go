package auth_cognito

import (
	"context"
	"monitoring-system/server/domain/auth"
	"monitoring-system/server/pkg/app_error"
	"monitoring-system/server/pkg/jwt_verify"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type cognitoAuth struct {
	client    *cognito.Client
	clientId  string
	ctx       context.Context
	jwtVerify jwt_verify.JWTVerify
}

func NewCognitoAuth(ctx context.Context, cognito *cognito.Client, clientId string, jwtVerify jwt_verify.JWTVerify) auth.CognitoAuth {
	return &cognitoAuth{
		client:    cognito,
		clientId:  clientId,
		ctx:       ctx,
		jwtVerify: jwtVerify,
	}
}

func (c *cognitoAuth) Login(input auth.LoginInput) (*auth.LoginOutput, error) {
	initiateAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": input.Username,
			"PASSWORD": input.Password,
		},
		ClientId: aws.String(c.clientId),
	}
	cognitoOut, err := c.client.InitiateAuth(c.ctx, initiateAuthInput)
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
	out := &auth.LoginOutput{
		AccessToken:  *cognitoOut.AuthenticationResult.AccessToken,
		RefreshToken: *cognitoOut.AuthenticationResult.RefreshToken,
		IdToken:      *cognitoOut.AuthenticationResult.IdToken,
	}

	return out, nil
}

func (c *cognitoAuth) SignUp(input auth.SignUpInput) (*auth.SignUpOutput, error) {
	singUpInput := &cognito.SignUpInput{
		ClientId: aws.String(c.clientId),
		Username: aws.String(input.Username),
		Password: aws.String(input.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(input.Name),
			},
		},
	}
	cognitoOut, err := c.client.SignUp(c.ctx, singUpInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UsernameExistsException") {
			return nil, app_error.NewApiError(409, "Username already exists")
		}
		return nil, err
	}

	out := &auth.SignUpOutput{
		IsConfirmed: cognitoOut.UserConfirmed,
	}
	return out, nil
}

// func (c *authServiceImpl) UserInformation(accessToken string) (*cognito.GetUserOutput, error) {
// 	input := &cognito.GetUserInput{
// 		AccessToken: aws.String(accessToken),
// 	}
// 	return c.client.GetUser(c.ctx, input)
// }

func (c *cognitoAuth) ConfirmSignUp(input auth.ConfirmSignUpInput) (*auth.ConfirmSignUpOutput, error) {
	confirmSignUp := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientId),
		Username:         aws.String(input.Username),
		ConfirmationCode: aws.String(input.Code),
	}

	_, err := c.client.ConfirmSignUp(c.ctx, confirmSignUp)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "CodeMismatchException") {
			return nil, app_error.NewApiError(400, "Invalid confirmation code")
		}
		if strings.Contains(errorType, "ExpiredCodeException") {
			return nil, app_error.NewApiError(400, "Confirmation code expired")
		}
		return nil, err
	}

	return &auth.ConfirmSignUpOutput{}, nil

}

func (c *cognitoAuth) GetUser(input auth.GetUserInput) (*auth.GetUserOutput, error) {
	getUserInput := &cognito.GetUserInput{
		AccessToken: &input.AccessToken,
	}
	cognitoOut, err := c.client.GetUser(c.ctx, getUserInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return nil, app_error.NewApiError(401, "Invalid access token")
		}
		if strings.Contains(errorType, "UserNotFoundException") {
			return nil, app_error.NewApiError(404, "User not found")
		}

		return nil, err
	}

	out := &auth.GetUserOutput{
		Username: *cognitoOut.Username,
		Name:     *cognitoOut.UserAttributes[0].Value, //TODO: handle this better
	}

	return out, nil

}

func (c *cognitoAuth) ValidateToken(token string) (*auth.Claims, error) {
	// err := c.jwtVerify.CacheJWK() //Check if this is the right place to cache the JWK
	// if err != nil {
	// 	return nil, err
	// }

	_, claims, err := c.jwtVerify.ParseJWT(token)
	if err != nil {
		return nil, err
	}

	return &auth.Claims{
		Email: claims.Email,
		Id:    claims.Sub,
	}, nil
}
