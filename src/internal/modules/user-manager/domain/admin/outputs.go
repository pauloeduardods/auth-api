package admin

import "context"

type UpdateAdminOutput struct {
	backup *Admin
	svc    AdminService
}

func NewUpdateAdminOutput(backup *Admin, svc AdminService) *UpdateAdminOutput {
	return &UpdateAdminOutput{
		backup: backup,
		svc:    svc,
	}
}

func (u *UpdateAdminOutput) Rollback(ctx context.Context) error {
	if u.backup == nil {
		return nil
	}
	updateAdminInput := &UpdateAdminInput{
		ID:    u.backup.ID,
		Email: &u.backup.Email,
		Name:  &u.backup.Name,
	}
	if err := updateAdminInput.Validate(); err != nil {
		return err
	}

	_, err := u.svc.Update(updateAdminInput)
	if err != nil {
		return err
	}
	return nil
}

type CreateAdminOutput struct {
	ID  *AdminID
	svc AdminService
}

func NewCreateAdminOutput(id *AdminID, svc AdminService) *CreateAdminOutput {
	return &CreateAdminOutput{
		ID:  id,
		svc: svc,
	}
}

func (c *CreateAdminOutput) Rollback(ctx context.Context) error {
	if c.ID == nil {
		return nil
	}

	deleteAdminInput := &DeleteAdminInput{
		ID: *c.ID,
	}
	if err := deleteAdminInput.Validate(); err != nil {
		return err
	}

	_, err := c.svc.Delete(deleteAdminInput)
	if err != nil {
		return err
	}
	return nil
}

type DeleteAdminOutput struct {
	backup *Admin
	repo   AdminRepository
}

func NewDeleteAdminOutput(backup *Admin, repo AdminRepository) *DeleteAdminOutput {
	return &DeleteAdminOutput{
		backup: backup,
		repo:   repo,
	}
}

func (d *DeleteAdminOutput) Rollback(ctx context.Context) error {
	if d.backup == nil {
		return nil
	}
	createAdminInput := &CreateAdminInput{
		Email: d.backup.Email,
		Name:  d.backup.Name,
	}
	if err := createAdminInput.Validate(); err != nil {
		return err
	}

	err := d.repo.Create(createAdminInput)
	if err != nil {
		return err
	}
	return nil
}
