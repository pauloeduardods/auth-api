package user

import "context"

type UpdateUserOutput struct {
	backup *User
	svc    UserService
}

func NewUpdateUserOutput(backup *User, svc UserService) *UpdateUserOutput {
	return &UpdateUserOutput{
		backup: backup,
		svc:    svc,
	}
}

func (u *UpdateUserOutput) Rollback(ctx context.Context) error {
	if u.backup == nil {
		return nil
	}
	updateUserInput := &UpdateUserInput{
		ID:    u.backup.ID,
		Email: &u.backup.Email,
		Name:  &u.backup.Name,
		Phone: u.backup.Phone,
	}
	if err := updateUserInput.Validate(); err != nil {
		return err
	}

	_, err := u.svc.Update(updateUserInput)
	if err != nil {
		return err
	}
	return nil
}

type CreateUserOutput struct {
	ID  *UserID
	svc UserService
}

func NewCreateUserOutput(id *UserID, svc UserService) *CreateUserOutput {
	return &CreateUserOutput{
		ID:  id,
		svc: svc,
	}
}

func (c *CreateUserOutput) Rollback(ctx context.Context) error {
	if c.ID == nil {
		return nil
	}

	deleteUserInput := &DeleteUserInput{
		ID: *c.ID,
	}
	if err := deleteUserInput.Validate(); err != nil {
		return err
	}

	_, err := c.svc.Delete(deleteUserInput)
	if err != nil {
		return err
	}
	return nil
}

type DeleteUserOutput struct {
	backup *User
	repo   UserRepository
}

func NewDeleteUserOutput(backup *User, repo UserRepository) *DeleteUserOutput {
	return &DeleteUserOutput{
		backup: backup,
		repo:   repo,
	}
}

func (d *DeleteUserOutput) Rollback(ctx context.Context) error {
	if d.backup == nil {
		return nil
	}
	createUserInput := &CreateUserInput{
		Email: d.backup.Email,
		Name:  d.backup.Name,
		Phone: d.backup.Phone,
	}
	if err := createUserInput.Validate(); err != nil {
		return err
	}

	err := d.repo.Create(createUserInput)
	if err != nil {
		return err
	}
	return nil
}
