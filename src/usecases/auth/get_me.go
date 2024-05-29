package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
)

type GetMeUseCase struct {
	auth auth.Auth
}

type GetMeInput struct {
	auth.GetMeInput
}

func NewGetMeUseCase(auth auth.Auth) *GetMeUseCase {
	return &GetMeUseCase{
		auth: auth,
	}
}

func (uc *GetMeUseCase) Execute(ctx context.Context, input GetMeInput) (*auth.GetMeOutput, error) {
	getMesInput, err := auth.NewGetMeInput(input.AccessToken)
	if err != nil {
		return nil, err
	}

	return uc.auth.GetMe(ctx, getMesInput)
}
