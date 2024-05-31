package auth_usecases

import (
	"auth-api/src/internal/domain/auth"
	"context"
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
	if err := input.RefreshTokenInput.Validate(); err != nil {
		return nil, err
	}

	return uc.auth.RefreshToken(ctx, input.RefreshTokenInput)
}
