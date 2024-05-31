package auth_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
)

type LoginUseCase struct {
	auth auth.AuthService
}

type LoginInput struct {
	auth.LoginInput
}

func NewLoginUseCase(auth auth.AuthService) *LoginUseCase {
	return &LoginUseCase{
		auth: auth,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*auth.LoginOutput, error) {
	loginInput, err := auth.NewLoginInput(input.Username, input.Password)
	if err != nil {
		return nil, err
	}

	return uc.auth.Login(ctx, loginInput)
}
