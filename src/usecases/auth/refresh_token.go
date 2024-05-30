package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
)

type RefreshTokenUseCase struct {
	auth auth.AuthService
}

type RefreshTokenInput struct {
	auth.RefreshTokenInput
}

func NewRefreshTokenUseCase(auth auth.AuthService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		auth: auth,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*auth.RefreshTokenOutput, error) {
	refreshTokenInput, err := auth.NewRefreshTokenInput(input.RefreshToken)
	if err != nil {
		return nil, err
	}

	return uc.auth.RefreshToken(ctx, refreshTokenInput)
}
