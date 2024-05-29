package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
)

type AddGroupUseCase struct {
	auth auth.Auth
}

type AddGroupInput struct {
	auth.AddGroupInput
}

func NewAddGroupUseCase(auth auth.Auth) *AddGroupUseCase {
	return &AddGroupUseCase{
		auth: auth,
	}
}

func (uc *AddGroupUseCase) Execute(ctx context.Context, input AddGroupInput) error {
	addGroupInput, err := auth.NewAddGroupInput(input.Username, input.GroupName)
	if err != nil {
		return err
	}

	return uc.auth.AddGroup(ctx, addGroupInput)
}
