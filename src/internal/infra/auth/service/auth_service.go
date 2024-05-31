package auth_service

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
	"monitoring-system/server/src/pkg/app_error"
	"monitoring-system/server/src/pkg/jwt_verify"
	"monitoring-system/server/src/pkg/logger"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type cognitoClient struct {
	client     *cognito.Client
	clientId   string
	userPoolId string
	jwtVerify  jwt_verify.JWTVerify
	logger     logger.Logger
}

func NewAuthService(cognito *cognito.Client, clientId string, jwtVerify jwt_verify.JWTVerify, userPoolId string, logger logger.Logger) auth.AuthService {
	return &cognitoClient{
		client:     cognito,
		clientId:   clientId,
		jwtVerify:  jwtVerify,
		userPoolId: userPoolId,
		logger:     logger,
	}
}

func (c *cognitoClient) AddMFA(ctx context.Context, input auth.AddMFAInput) (*auth.AddMFAOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	associateSoftwareTokenInput := &cognito.AssociateSoftwareTokenInput{
		AccessToken: aws.String(input.AccessToken),
	}

	associateSoftwareTokenOutput, err := c.client.AssociateSoftwareToken(ctx, associateSoftwareTokenInput)
	if err != nil {
		c.logger.Error("Cognito associate software token error", err)
		return nil, app_error.NewApiError(500, "Failed to associate software token")
	}

	return &auth.AddMFAOutput{
		SecretCode: *associateSoftwareTokenOutput.SecretCode,
		Session:    associateSoftwareTokenOutput.Session,
	}, nil
}

func (c *cognitoClient) ActivateMFA(ctx context.Context, input auth.ActivateMFAInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifySoftwareTokenInput := &cognito.VerifySoftwareTokenInput{
		UserCode:    aws.String(input.Code),
		AccessToken: aws.String(input.AccessToken),
	}

	verifyOut, err := c.client.VerifySoftwareToken(ctx, verifySoftwareTokenInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "CodeMismatchException") {
			return app_error.NewApiError(400, "Invalid MFA code")
		}
		if strings.Contains(errorType, "NotAuthorizedException") {
			return app_error.NewApiError(401, "Invalid access token")
		}
		c.logger.Error("Cognito verify software token error", err)
		return app_error.NewApiError(400, "Failed to verify software token")
	}
	if verifyOut.Status != "SUCCESS" {
		return app_error.NewApiError(400, "Invalid MFA code")
	}

	setUserMFAPreferenceInput := &cognito.SetUserMFAPreferenceInput{
		AccessToken: aws.String(input.AccessToken),
		SoftwareTokenMfaSettings: &types.SoftwareTokenMfaSettingsType{
			Enabled:      true,
			PreferredMfa: true,
		},
	}

	_, err = c.client.SetUserMFAPreference(ctx, setUserMFAPreferenceInput)
	if err != nil {
		c.logger.Error("Cognito set user MFA preference error", err)
		return app_error.NewApiError(500, "Failed to set user MFA preference")
	}

	return nil
}

func (c *cognitoClient) VerifyMFA(ctx context.Context, input auth.VerifyMFAInput) (*auth.LoginOutput, error) {
	respondToAuthChallengeInput := &cognito.RespondToAuthChallengeInput{
		ChallengeName: "SOFTWARE_TOKEN_MFA",
		ClientId:      aws.String(c.clientId),
		Session:       aws.String(input.Session),
		ChallengeResponses: map[string]string{
			"USERNAME":                input.Username,
			"SOFTWARE_TOKEN_MFA_CODE": input.Code,
		},
	}

	cognitoOut, err := c.client.RespondToAuthChallenge(ctx, respondToAuthChallengeInput)
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

func (c *cognitoClient) AdminRemoveMFA(ctx context.Context, input auth.AdminRemoveMFAInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	adminSetUserMFAPreferenceInput := &cognito.AdminSetUserMFAPreferenceInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		SoftwareTokenMfaSettings: &types.SoftwareTokenMfaSettingsType{
			Enabled:      false,
			PreferredMfa: false,
		},
	}

	_, err := c.client.AdminSetUserMFAPreference(ctx, adminSetUserMFAPreferenceInput)
	if err != nil {
		c.logger.Error("Cognito remove MFA error", err)
		return app_error.NewApiError(500, "Failed to remove MFA")
	}

	return nil
}

func (c *cognitoClient) RemoveMFA(ctx context.Context, input auth.RemoveMFAInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	setUserMFAPreferenceInput := &cognito.SetUserMFAPreferenceInput{
		AccessToken: aws.String(input.AccessToken),
		SoftwareTokenMfaSettings: &types.SoftwareTokenMfaSettingsType{
			Enabled:      false,
			PreferredMfa: false,
		},
	}

	_, err := c.client.SetUserMFAPreference(ctx, setUserMFAPreferenceInput)
	if err != nil {
		c.logger.Error("Cognito remove MFA error", err)
		return app_error.NewApiError(500, "Failed to remove MFA")
	}

	return nil
}

func (c *cognitoClient) Login(ctx context.Context, input auth.LoginInput) (*auth.LoginOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	initiateAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": input.Username,
			"PASSWORD": input.Password,
		},
		ClientId: aws.String(c.clientId),
	}
	cognitoOut, err := c.client.InitiateAuth(ctx, initiateAuthInput)
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

func (c *cognitoClient) SignUp(ctx context.Context, input auth.SignUpInput) (*auth.SignUpOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

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
	cognitoOut, err := c.client.SignUp(ctx, signUpInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UsernameExistsException") {
			return nil, app_error.NewApiError(409, "Username already exists")
		}
		c.logger.Error("Cognito signup error", err)
		return nil, err
	}

	err = c.AddGroup(ctx, auth.AddGroupInput{
		Username:  input.Username,
		GroupName: auth.User,
	})
	if err != nil {
		return nil, err //TODO: Rollback signup
	}

	out := &auth.SignUpOutput{
		IsConfirmed: cognitoOut.UserConfirmed,
		Id:          *cognitoOut.UserSub,
	}
	return out, nil
}

func (c *cognitoClient) ConfirmSignUp(ctx context.Context, input auth.ConfirmSignUpInput) (*auth.ConfirmSignUpOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	confirmSignUp := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientId),
		Username:         aws.String(input.Username),
		ConfirmationCode: aws.String(input.Code),
	}

	_, err := c.client.ConfirmSignUp(ctx, confirmSignUp)
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

func (c *cognitoClient) GetMe(ctx context.Context, input auth.GetMeInput) (*auth.GetMeOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	getMeInput := &cognito.GetUserInput{
		AccessToken: &input.AccessToken,
	}
	cognitoOut, err := c.client.GetUser(ctx, getMeInput)
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

	out := &auth.GetMeOutput{
		Username: *cognitoOut.Username,
		Name:     *cognitoOut.UserAttributes[0].Value, //TODO: handle this better
	}

	return out, nil
}

func (c *cognitoClient) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
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

func (c *cognitoClient) AddGroup(ctx context.Context, input auth.AddGroupInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	addUserToGroupInput := &cognito.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		GroupName:  aws.String(string(input.GroupName)),
	}

	_, err := c.client.AdminAddUserToGroup(ctx, addUserToGroupInput)
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

func (c *cognitoClient) RemoveGroup(ctx context.Context, input auth.RemoveGroupInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	removeUserFromGroupInput := &cognito.AdminRemoveUserFromGroupInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		GroupName:  aws.String(string(input.GroupName)),
	}

	_, err := c.client.AdminRemoveUserFromGroup(ctx, removeUserFromGroupInput)
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

func (c *cognitoClient) RefreshToken(ctx context.Context, input auth.RefreshTokenInput) (*auth.RefreshTokenOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	refreshTokenInput := &cognito.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": input.RefreshToken,
		},
		ClientId: aws.String(c.clientId),
	}
	cognitoOut, err := c.client.InitiateAuth(ctx, refreshTokenInput)
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

func (c *cognitoClient) CreateAdmin(ctx context.Context, input auth.CreateAdminInput) (*auth.CreateAdminOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

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

	_, err := c.client.AdminCreateUser(ctx, createUserInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UsernameExistsException") {
			return nil, app_error.NewApiError(409, "Username already exists")
		}
		c.logger.Error("Cognito admin create user error", err)
		return nil, err
	}

	err = c.AddGroup(ctx, auth.AddGroupInput{
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

func (c *cognitoClient) DeleteUser(ctx context.Context, input auth.DeleteUserInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	deleteUserInput := &cognito.AdminDeleteUserInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
	}

	_, err := c.client.AdminDeleteUser(ctx, deleteUserInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return app_error.NewApiError(404, "User not found")
		}
		c.logger.Error("Cognito delete user error", err)
		return err
	}

	return nil
}

func (c *cognitoClient) Logout(ctx context.Context, input auth.LogoutInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	globalSignOutInput := &cognito.GlobalSignOutInput{
		AccessToken: aws.String(input.AccessToken),
	}

	_, err := c.client.GlobalSignOut(ctx, globalSignOutInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return app_error.NewApiError(401, "Invalid access token")
		}
		c.logger.Error("Cognito logout error", err)
		return err
	}

	return nil
}
