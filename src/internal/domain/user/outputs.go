package user

import "context"

type UpdateUserOutput struct {
	Backup *User
}

func (u *UpdateUserOutput) Rollback(ctx context.Context, userService UserService) error {
	if u.Backup == nil {
		return nil
	}
	updateUserInput := &UpdateUserInput{
		ID:    u.Backup.ID,
		Email: &u.Backup.Email,
		Name:  &u.Backup.Name,
		Phone: u.Backup.Phone,
	}
	if err := updateUserInput.Validate(); err != nil {
		return err
	}

	_, err := userService.Update(updateUserInput)
	if err != nil {
		return err
	}
	return nil
}

type CreateUserOutput struct {
	ID *UserID
}

func (c *CreateUserOutput) Rollback(ctx context.Context, userService UserService) error {
	if c.ID == nil {
		return nil
	}

	deleteUserInput := &DeleteUserInput{
		ID: *c.ID,
	}
	if err := deleteUserInput.Validate(); err != nil {
		return err
	}

	_, err := userService.Delete(deleteUserInput)
	if err != nil {
		return err
	}
	return nil
}

type DeleteUserOutput struct {
	Backup *User
}

func (d *DeleteUserOutput) Rollback(ctx context.Context, authRepo UserRepository) error {
	if d.Backup == nil {
		return nil
	}
	createUserInput := &CreateUserInput{
		Email: d.Backup.Email,
		Name:  d.Backup.Name,
		Phone: d.Backup.Phone,
	}
	if err := createUserInput.Validate(); err != nil {
		return err
	}

	err := authRepo.Create(createUserInput)
	if err != nil {
		return err
	}
	return nil
}
