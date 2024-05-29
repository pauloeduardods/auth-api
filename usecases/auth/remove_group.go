package auth

import (
	"context"
	"monitoring-system/server/domain/auth"
)

type RemoveGroupUseCase struct {
	auth auth.Auth
}

type RemoveGroupInput struct {
	auth.RemoveGroupInput
}

func NewRemoveGroupUseCase(auth auth.Auth) *RemoveGroupUseCase {
	return &RemoveGroupUseCase{
		auth: auth,
	}
}

func (uc *RemoveGroupUseCase) Execute(ctx context.Context, input RemoveGroupInput) error {
	removeGroupInput, err := auth.NewRemoveGroupInput(input.Username, input.GroupName)
	if err != nil {
		return err
	}

	return uc.auth.RemoveGroup(ctx, removeGroupInput)
}
