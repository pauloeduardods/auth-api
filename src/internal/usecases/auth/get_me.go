package auth_usecases

import (
	"auth-api/src/internal/domain/auth"
	"context"
)

type GetMeUseCase struct {
	auth auth.AuthService
}

type GetMeInput struct {
	auth.GetMeInput
}

func NewGetMeUseCase(auth auth.AuthService) *GetMeUseCase {
	return &GetMeUseCase{
		auth: auth,
	}
}

func (uc *GetMeUseCase) Execute(ctx context.Context, input GetMeInput) (*auth.GetMeOutput, error) {
	if err := input.GetMeInput.Validate(); err != nil {
		return nil, err
	}

	return uc.auth.GetMe(ctx, input.GetMeInput)
}
