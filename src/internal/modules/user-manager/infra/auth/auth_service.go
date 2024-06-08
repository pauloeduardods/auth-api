package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/internal/modules/user-manager/domain/user"
	"auth-api/src/internal/shared/code/domain/code"
	"auth-api/src/internal/shared/notification/domain/email"
	"auth-api/src/pkg/app_error"
	"auth-api/src/pkg/jwt_verify"
	"auth-api/src/pkg/logger"
	"context"
	"fmt"
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
	email      email.EmailService
	code       code.CodeService
}

func NewAuthService(cognito *cognito.Client, clientId string, jwtVerify jwt_verify.JWTVerify, userPoolId string, logger logger.Logger, email email.EmailService, code code.CodeService) auth.AuthService {
	return &cognitoClient{
		client:     cognito,
		clientId:   clientId,
		jwtVerify:  jwtVerify,
		userPoolId: userPoolId,
		logger:     logger,
		email:      email,
		code:       code,
	}
}

func (c *cognitoClient) AddMFA(ctx context.Context, input auth.AddMFAInput) (*auth.AddMFAOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

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
	if err := input.Validate(); err != nil {
		return err
	}

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
			return auth.ErrInvalidMfaCode
		}
		if strings.Contains(errorType, "NotAuthorizedException") {
			return auth.ErrInvalidAccessCode
		}
		c.logger.Error("Cognito verify software token error", err)
		return auth.ErrFailedToVerifySoftwareMfa
	}
	if verifyOut.Status != "SUCCESS" {
		return auth.ErrInvalidMfaCode
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
	if err := input.Validate(); err != nil {
		return nil, err
	}

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
		return nil, auth.ErrFailedToRespondToChallenge
	}

	return &auth.LoginOutput{
		AccessToken:  cognitoOut.AuthenticationResult.AccessToken,
		RefreshToken: cognitoOut.AuthenticationResult.RefreshToken,
		IdToken:      cognitoOut.AuthenticationResult.IdToken,
	}, nil
}

func (c *cognitoClient) AdminRemoveMFA(ctx context.Context, input auth.AdminRemoveMFAInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

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
	if err := input.Validate(); err != nil {
		return err
	}

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
	if err := input.Validate(); err != nil {
		return nil, err
	}

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
			return nil, auth.ErrInvalidUsernameOrPassword
		}
		if strings.Contains(errorType, "PasswordResetRequiredException") {
			return nil, auth.ErrPasswordResetRequired
		}
		if strings.Contains(errorType, "UserNotConfirmedException") {
			return nil, auth.ErrUserNotConfirmed
		}
		c.logger.Error("Cognito login error", err)
		return nil, err
	}

	if cognitoOut.ChallengeName != "" {
		return &auth.LoginOutput{
			Session:  cognitoOut.Session,
			NextStep: (*string)(&cognitoOut.ChallengeName),
		}, nil
	}

	if cognitoOut.AuthenticationResult == nil {
		return nil, app_error.NewApiError(500, "Failed to get authentication result")
	}

	out := &auth.LoginOutput{
		AccessToken:  cognitoOut.AuthenticationResult.AccessToken,
		RefreshToken: cognitoOut.AuthenticationResult.RefreshToken,
		IdToken:      cognitoOut.AuthenticationResult.IdToken,
	}

	return out, nil
}

func (c *cognitoClient) SignUp(ctx context.Context, input auth.SignUpInput) (o *auth.SignUpOutput, execErr error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

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
			return nil, auth.ErrUserAlreadyExists
		}
		c.logger.Error("Cognito signup error", err)
		return nil, err
	}

	out := auth.NewSignUpOutput(*cognitoOut.UserSub, input.Username, false, c)

	defer func() {
		if execErr != nil {
			c.logger.Info("Rollback signup")
			if err := out.Rollback(ctx); err != nil {
				c.logger.Error("Rollback signup error", err)
			}
		}
	}()

	err = c.AddGroup(ctx, auth.AddGroupInput{
		Username:  input.Username,
		GroupName: auth.User,
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *cognitoClient) ConfirmSignUp(ctx context.Context, input auth.ConfirmSignUpInput) (*auth.ConfirmSignUpOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	adminConfirmSignUpInput := &cognito.AdminConfirmSignUpInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
	}

	_, err := c.client.AdminConfirmSignUp(ctx, adminConfirmSignUpInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return nil, auth.ErrUserNotFound
		}
		c.logger.Error("Cognito confirm signup error", err)
		return nil, err
	}

	return &auth.ConfirmSignUpOutput{}, nil
}

func (c *cognitoClient) GetMe(ctx context.Context, input auth.GetMeInput) (*auth.GetMeOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	getMeInput := &cognito.GetUserInput{
		AccessToken: &input.AccessToken,
	}
	cognitoOut, err := c.client.GetUser(ctx, getMeInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return nil, auth.ErrInvalidAccessCode
		}
		if strings.Contains(errorType, "UserNotFoundException") {
			return nil, user.ErrUserNotFound
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
			return auth.ErrUserNotFound
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
	if err := input.Validate(); err != nil {
		return err
	}

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
			return auth.ErrUserNotFound
		}
		if strings.Contains(errorType, "ResourceNotFoundException") {
			return auth.ErrInvalidGroup
		}
		c.logger.Error("Cognito remove group error", err)
		return err
	}

	return nil
}

func (c *cognitoClient) RefreshToken(ctx context.Context, input auth.RefreshTokenInput) (*auth.RefreshTokenOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

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
			return nil, auth.ErrInvalidRefreshToken
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

func (c *cognitoClient) CreateAdmin(ctx context.Context, input auth.CreateAdminInput) (o *auth.CreateAdminOutput, execErr error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

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

	cognitoOut, err := c.client.AdminCreateUser(ctx, createUserInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UsernameExistsException") {
			return nil, auth.ErrUserAlreadyExists
		}
		c.logger.Error("Cognito admin create user error", err)
		return nil, err
	}

	var userId string
	for _, attr := range cognitoOut.User.Attributes {
		if *attr.Name == "sub" {
			userId = *attr.Value
			break
		}
	}

	out := auth.NewCreateAdminOutput(userId, input.Username, c)
	defer func() {
		if execErr != nil {
			if err := out.Rollback(ctx); err != nil {
				c.logger.Error("Rollback create admin error", err)
			}
		}
	}()

	if userId == "" {
		return nil, app_error.NewApiError(500, "Failed to get user id")
	}

	err = c.AddGroup(ctx, auth.AddGroupInput{
		Username:  input.Username,
		GroupName: auth.Admin,
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *cognitoClient) DeleteUser(ctx context.Context, input auth.DeleteUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

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
			return auth.ErrUserNotFound
		}
		c.logger.Error("Cognito delete user error", err)
		return err
	}

	return nil
}

func (c *cognitoClient) Logout(ctx context.Context, input auth.LogoutInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	globalSignOutInput := &cognito.GlobalSignOutInput{
		AccessToken: aws.String(input.AccessToken),
	}

	_, err := c.client.GlobalSignOut(ctx, globalSignOutInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return auth.ErrInvalidAccessCode
		}
		c.logger.Error("Cognito logout error", err)
		return err
	}

	return nil
}

func (c *cognitoClient) SetPassword(ctx context.Context, input auth.SetPasswordInput) (o *auth.LoginOutput, execErr error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	respondToAuthChallengeInput := &cognito.RespondToAuthChallengeInput{
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		ClientId:      aws.String(c.clientId),
		ChallengeResponses: map[string]string{
			"USERNAME":     input.Username,
			"NEW_PASSWORD": input.Password,
		},
		Session: aws.String(input.Session),
	}

	authOut, err := c.client.RespondToAuthChallenge(ctx, respondToAuthChallengeInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return nil, auth.ErrUserNotFound
		}
		c.logger.Error("Cognito set password error", err)
		return nil, err
	}

	return &auth.LoginOutput{
		AccessToken:  authOut.AuthenticationResult.AccessToken,
		RefreshToken: authOut.AuthenticationResult.RefreshToken,
		IdToken:      authOut.AuthenticationResult.IdToken,
	}, nil
}

func (c *cognitoClient) GetUser(ctx context.Context, input auth.GetUserInput) (*auth.GetUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	getUserInput := &cognito.AdminGetUserInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
	}

	cognitoOut, err := c.client.AdminGetUser(ctx, getUserInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return nil, auth.ErrUserNotFound
		}
		c.logger.Error("Cognito get user error", err)
		return nil, err
	}

	var username, name, id string

	for _, attr := range cognitoOut.UserAttributes {
		switch *attr.Name {
		case "email":
			username = *attr.Value
		case "name":
			name = *attr.Value
		case "sub":
			id = *attr.Value
		}
	}

	out := &auth.GetUserOutput{
		Username: username,
		Name:     name,
		Id:       id,
	}

	return out, nil
}

func (c *cognitoClient) AdminLogout(ctx context.Context, input auth.AdminLogoutInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	adminUserGlobalSignOutInput := &cognito.AdminUserGlobalSignOutInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
	}

	_, err := c.client.AdminUserGlobalSignOut(ctx, adminUserGlobalSignOutInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return auth.ErrUserNotFound
		}
		c.logger.Error("Cognito admin logout error", err)
		return err
	}

	return nil
}

func (c *cognitoClient) VerifyEmail(ctx context.Context, input auth.VerifyEmailInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyUserAttributeInput := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
		},
	}

	_, err := c.client.AdminUpdateUserAttributes(ctx, verifyUserAttributeInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return auth.ErrUserNotFound
		}
		c.logger.Error("Cognito verify email error", err)
		return err
	}

	return nil
}

func (c *cognitoClient) GenerateAndSendCode(ctx context.Context, input auth.GenerateAndSendCodeInput) (*auth.GenerateAndSendCodeOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	expiresAt := time.Now().Add(10 * time.Minute)

	generateAndSendInput := &code.GenerateAndSaveInput{
		Identifier:        fmt.Sprintf("%s#%s", input.Identifier, input.Username),
		ExpiresAt:         expiresAt,
		Length:            6,
		CanContainLetters: false,
	}

	code, err := c.code.GenerateAndSave(ctx, *generateAndSendInput)
	if err != nil {
		return nil, err
	}

	sendEmailInput := email.Email{
		To:      input.Username,
		Subject: input.Subject,
		Body:    fmt.Sprintf(input.Body, code.Value),
	}
	if err := c.email.SendEmail(ctx, sendEmailInput); err != nil {
		return nil, err
	}

	return &auth.GenerateAndSendCodeOutput{
		Code: code.Value,
	}, nil

}

func (c *cognitoClient) VerifyCode(ctx context.Context, input auth.VerifyCodeInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	identifier := fmt.Sprintf("%s#%s", input.Identifier, input.Username)

	verifyInput := code.VerifyCodeInput{
		Identifier: identifier,
		Code:       input.Code,
	}

	if err := c.code.VerifyCode(ctx, verifyInput); err != nil {
		return err
	}

	return nil
}

func (c *cognitoClient) ChangeForgotPassword(ctx context.Context, input auth.ChangeForgotPasswordInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	admSetPassword := &cognito.AdminSetUserPasswordInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(input.Username),
		Password:   aws.String(input.NewPassword),
		Permanent:  true,
	}

	_, err := c.client.AdminSetUserPassword(ctx, admSetPassword)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "UserNotFoundException") {
			return auth.ErrUserNotFound
		}
		c.logger.Error("Cognito change forgot password error", err)
		return err
	}

	return nil
}

func (c *cognitoClient) ChangePassword(ctx context.Context, input auth.ChangePasswordInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	changePasswordInput := &cognito.ChangePasswordInput{
		AccessToken:      aws.String(input.AccessToken),
		PreviousPassword: aws.String(input.OldPassword),
		ProposedPassword: aws.String(input.NewPassword),
	}

	_, err := c.client.ChangePassword(ctx, changePasswordInput)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return auth.ErrInvalidAccessCode
		}
		c.logger.Error("Cognito change password error", err)
		return err
	}

	return nil
}
