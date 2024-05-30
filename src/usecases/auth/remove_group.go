package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
)

type RemoveGroupUseCase struct {
	auth auth.AuthService
}

type RemoveGroupInput struct {
	auth.RemoveGroupInput
}

func NewRemoveGroupUseCase(auth auth.AuthService) *RemoveGroupUseCase {
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
