package auth

import (
	"context"
	"monitoring-system/server/domain/auth"
)

type VerifyMFAUseCase struct {
	auth auth.Auth
}

type VerifyMFAInput struct {
	auth.VerifyMFAInput
}

func NewVerifyMFAUseCase(auth auth.Auth) *VerifyMFAUseCase {
	return &VerifyMFAUseCase{
		auth: auth,
	}
}

func (uc *VerifyMFAUseCase) Execute(ctx context.Context, input VerifyMFAInput) (*auth.LoginOutput, error) {
	verifyMFAInput, err := auth.NewVerifyMFAInput(input.Code, input.Username, input.Session)
	if err != nil {
		return nil, err
	}

	return uc.auth.VerifyMFA(ctx, verifyMFAInput)
}
