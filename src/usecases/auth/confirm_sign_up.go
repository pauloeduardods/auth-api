package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
)

type ConfirmSignUpUseCase struct {
	auth auth.Auth
}

type ConfirmSignUpInput struct {
	auth.ConfirmSignUpInput
}

func NewConfirmSignUpUseCase(auth auth.Auth) *ConfirmSignUpUseCase {
	return &ConfirmSignUpUseCase{
		auth: auth,
	}
}

func (uc *ConfirmSignUpUseCase) Execute(ctx context.Context, input ConfirmSignUpInput) (*auth.ConfirmSignUpOutput, error) {
	confirmSignUpInput, err := auth.NewConfirmSignUpInput(input.Username, input.Code)
	if err != nil {
		return nil, err
	}

	return uc.auth.ConfirmSignUp(ctx, confirmSignUpInput)
}
