package auth_cognito

import (
	"context"
	"monitoring-system/server/domain/auth"
	"monitoring-system/server/pkg/app_error"
	"monitoring-system/server/pkg/jwt_verify"
	"monitoring-system/server/pkg/logger"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type cognitoAuth struct {
	client     *cognito.Client
	clientId   string
	userPoolId string
	ctx        context.Context
	jwtVerify  jwt_verify.JWTVerify
	logger     logger.Logger
}

func NewCognitoAuth(ctx context.Context, cognito *cognito.Client, clientId string, jwtVerify jwt_verify.JWTVerify, userPoolId string, logger logger.Logger) auth.CognitoAuth {
	return &cognitoAuth{
		client:     cognito,
		clientId:   clientId,
		ctx:        ctx,
		jwtVerify:  jwtVerify,
		userPoolId: userPoolId,
		logger:     logger,
	}
}

func (c *cognitoAuth) AddMFA(input auth.AddMFAInput) (*auth.AddMFAOutput, error) {
	associateSoftwareTokenInput := &cognito.AssociateSoftwareTokenInput{
		AccessToken: aws.String(input.AccessToken),
	}

	associateSoftwareTokenOutput, err := c.client.AssociateSoftwareToken(c.ctx, associateSoftwareTokenInput)
	if err != nil {
		c.logger.Error("Cognito associate software token error", err)
		return nil, app_error.NewApiError(500, "Failed to associate software token")
	}

	return &auth.AddMFAOutput{
		SecretCode: *associateSoftwareTokenOutput.SecretCode,
	}, nil
}

func (c *cognitoAuth) VerifyMFA(input auth.VerifyMFAInput) (*auth.LoginOutput, error) {
	// verifySoftwareTokenInput := &cognito.VerifySoftwareTokenInput{
	// 	UserCode:    aws.String(input.Code),
	// 	AccessToken: aws.String(input.AccessToken),
	// 	Session:     aws.String(input.Session),
	// }

	// verifySoftwareTokenOutput, err := c.client.VerifySoftwareToken(c.ctx, verifySoftwareTokenInput)
	// if err != nil {
	// 	c.logger.Error("Cognito verify software token error", err)
	// 	return nil, app_error.NewApiError(400, "Failed to verify software token")
	// }

	// if verifySoftwareTokenOutput.Status != "SUCCESS" {
	// 	return nil, app_error.NewApiError(400, "Invalid MFA code")
	// }

	respondToAuthChallengeInput := &cognito.RespondToAuthChallengeInput{
		ChallengeName: "SOFTWARE_TOKEN_MFA",
		ClientId:      aws.String(c.clientId),
		Session:       aws.String(input.Session),
		ChallengeResponses: map[string]string{
			"USERNAME":                input.Username,
			"SOFTWARE_TOKEN_MFA_CODE": input.Code,
		},
	}

	cognitoOut, err := c.client.RespondToAuthChallenge(c.ctx, respondToAuthChallengeInput)
	if err != nil {
		c.logger.Error("Cognito respond to auth challenge error", err)
		return nil, app_error.NewApiError(400, "Failed to respond to auth challenge")
	}

	return &auth.LoginOutput{
		AccessToken:  *cognitoOut.AuthenticationResult.AccessToken,
		RefreshToken: *cognitoOut.AuthenticationResult.RefreshToken,
		IdToken:      *cognitoOut.AuthenticationResult.IdToken,
	}, nil
}

func (c *cognitoAuth) RemoveMFA(input auth.RemoveMFAInput) error {
	adminSetUserMFAPreferenceInput := &cognito.AdminSetUserMFAPreferenceInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		SoftwareTokenMfaSettings: &types.SoftwareTokenMfaSettingsType{
			Enabled:      false,
			PreferredMfa: false,
		},
	}

	_, err := c.client.AdminSetUserMFAPreference(c.ctx, adminSetUserMFAPreferenceInput)
	if err != nil {
		c.logger.Error("Cognito remove MFA error", err)
		return app_error.NewApiError(500, "Failed to remove MFA")
	}

	return nil
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
		c.logger.Error("Cognito login error", err)
		return nil, err
	}

	if cognitoOut.ChallengeName == "SOFTWARE_TOKEN_MFA" {
		return &auth.LoginOutput{
			Session: *cognitoOut.Session,
		}, nil
	}

	out := &auth.LoginOutput{
		AccessToken:  *cognitoOut.AuthenticationResult.AccessToken,
		RefreshToken: *cognitoOut.AuthenticationResult.RefreshToken,
		IdToken:      *cognitoOut.AuthenticationResult.IdToken,
	}

	return out, nil
}

func (c *cognitoAuth) SignUp(input auth.SignUpInput) (*auth.SignUpOutput, error) {
	signUpInput := &cognito.SignUpInput{
		ClientId: aws.String(c.clientId),
		Username: aws.String(input.Username),
		Password: aws.String(input.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.Username),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(input.Name),
			},
		},
	}
	cognitoOut, err := c.client.SignUp(c.ctx, signUpInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UsernameExistsException") {
			return nil, app_error.NewApiError(409, "Username already exists")
		}
		c.logger.Error("Cognito signup error", err)
		return nil, err
	}

	err = c.AddGroup(auth.AddGroupInput{
		Username:  input.Username,
		GroupName: auth.User,
	})
	if err != nil {
		return nil, err //TODO: Rollback signup
	}

	out := &auth.SignUpOutput{
		IsConfirmed: cognitoOut.UserConfirmed,
	}
	return out, nil
}

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
		c.logger.Error("Cognito confirm signup error", err)
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
		c.logger.Error("Cognito get user error", err)
		return nil, err
	}

	out := &auth.GetUserOutput{
		Username: *cognitoOut.Username,
		Name:     *cognitoOut.UserAttributes[0].Value, //TODO: handle this better
	}

	return out, nil
}

func (c *cognitoAuth) ValidateToken(token string) (*auth.Claims, error) {
	_, claims, err := c.jwtVerify.ParseJWT(token)
	if err != nil {
		return nil, err
	}

	return &auth.Claims{
		Email:      claims.Email,
		Id:         claims.Sub,
		UserGroups: claims.UserGroups,
	}, nil
}

func (c *cognitoAuth) AddGroup(input auth.AddGroupInput) error {
	addUserToGroupInput := &cognito.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		GroupName:  aws.String(string(input.GroupName)),
	}

	_, err := c.client.AdminAddUserToGroup(c.ctx, addUserToGroupInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return app_error.NewApiError(404, "User not found")
		}
		if strings.Contains(errorType, "ResourceNotFoundException") {
			return app_error.NewApiError(404, "Group not found")
		}
		c.logger.Error("Cognito add group error", err)
		return err
	}

	return nil
}

func (c *cognitoAuth) RemoveGroup(input auth.RemoveGroupInput) error {
	removeUserFromGroupInput := &cognito.AdminRemoveUserFromGroupInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		GroupName:  aws.String(string(input.GroupName)),
	}

	_, err := c.client.AdminRemoveUserFromGroup(c.ctx, removeUserFromGroupInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return app_error.NewApiError(404, "User not found")
		}
		if strings.Contains(errorType, "ResourceNotFoundException") {
			return app_error.NewApiError(404, "Group not found")
		}
		c.logger.Error("Cognito remove group error", err)
		return err
	}

	return nil
}

func (c *cognitoAuth) RefreshToken(input auth.RefreshTokenInput) (*auth.RefreshTokenOutput, error) {
	refreshTokenInput := &cognito.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": input.RefreshToken,
		},
		ClientId: aws.String(c.clientId),
	}
	cognitoOut, err := c.client.InitiateAuth(c.ctx, refreshTokenInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return nil, app_error.NewApiError(401, "Invalid refresh token")
		}
		c.logger.Error("Cognito refresh token error", err)
		return nil, err
	}

	out := &auth.RefreshTokenOutput{
		AccessToken: *cognitoOut.AuthenticationResult.AccessToken,
		IdToken:     *cognitoOut.AuthenticationResult.IdToken,
	}

	return out, nil
}

func (c *cognitoAuth) CreateAdmin(input auth.CreateAdminInput) (*auth.CreateAdminOutput, error) {
	createUserInput := &cognito.AdminCreateUserInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.Username),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(input.Name),
			},
		},
		TemporaryPassword: aws.String(input.Password),
		DesiredDeliveryMediums: []types.DeliveryMediumType{
			types.DeliveryMediumTypeEmail,
		},
		ForceAliasCreation: true,
	}

	_, err := c.client.AdminCreateUser(c.ctx, createUserInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UsernameExistsException") {
			return nil, app_error.NewApiError(409, "Username already exists")
		}
		c.logger.Error("Cognito admin create user error", err)
		return nil, err
	}

	err = c.AddGroup(auth.AddGroupInput{
		Username:  input.Username,
		GroupName: auth.Admin,
	})
	if err != nil {
		return nil, err //TODO: Rollback create admin
	}

	return &auth.CreateAdminOutput{
		Username: input.Username,
	}, nil
}
