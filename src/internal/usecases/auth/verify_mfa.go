package auth_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
)

type VerifyMFAUseCase struct {
	auth auth.AuthService
}

type VerifyMFAInput struct {
	auth.VerifyMFAInput
}

func NewVerifyMFAUseCase(auth auth.AuthService) *VerifyMFAUseCase {
	return &VerifyMFAUseCase{
		auth: auth,
	}
}

func (uc *VerifyMFAUseCase) Execute(ctx context.Context, input VerifyMFAInput) (*auth.LoginOutput, error) {
	if err := input.VerifyMFAInput.Validate(); err != nil {
		return nil, err
	}

	return uc.auth.VerifyMFA(ctx, input.VerifyMFAInput)
}