package admin

import (
	"auth-api/src/internal/modules/user-manager/domain/admin"
	"auth-api/src/pkg/logger"
)

type AdminService struct {
	repo   admin.AdminRepository
	logger logger.Logger
}

func NewAdminService(repo admin.AdminRepository, logger logger.Logger) admin.AdminService {
	return &AdminService{
		repo:   repo,
		logger: logger,
	}
}

func (a *AdminService) GetByID(input *admin.GetAdminInput) (*admin.Admin, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	return a.repo.GetByID(input)
}

func (a *AdminService) GetByEmail(email *admin.GetAdminByEmailInput) (*admin.Admin, error) {
	if err := email.Validate(); err != nil {
		return nil, err
	}

	return a.repo.GetByEmail(email)
}

func (a *AdminService) Create(input *admin.CreateAdminInput) (*admin.CreateAdminOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	getAdminByEmailInput := admin.GetAdminByEmailInput{Email: input.Email}
	if err := getAdminByEmailInput.Validate(); err != nil {
		return nil, err
	}

	adminExists, err := a.repo.GetByEmail(&getAdminByEmailInput)
	if err != nil {
		if err != admin.ErrAdminNotFound {
			return nil, err
		}
	}

	if adminExists != nil {
		return nil, admin.ErrAdminAlreadyExists
	}

	out := admin.NewCreateAdminOutput(&input.ID, a)

	if err := a.repo.Create(input); err != nil {
		return nil, err
	}

	return out, nil
}

func (a *AdminService) Update(input *admin.UpdateAdminInput) (*admin.UpdateAdminOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	getAdminInput := admin.GetAdminInput{ID: input.ID.String()}
	if err := getAdminInput.Validate(); err != nil {
		return nil, err
	}

	adminOut, err := a.repo.GetByID(&getAdminInput)
	if err != nil {
		return nil, err
	}
	if adminOut == nil {
		return nil, admin.ErrAdminNotFound
	}

	out := admin.NewUpdateAdminOutput(adminOut, a)

	return out, a.repo.Update(input)
}

func (a *AdminService) Delete(id *admin.DeleteAdminInput) (*admin.DeleteAdminOutput, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}

	getAdminInput := admin.GetAdminInput{ID: id.ID.String()}
	if err := getAdminInput.Validate(); err != nil {
		return nil, err
	}

	adminOut, err := a.repo.GetByID(&getAdminInput)
	if err != nil {
		return nil, err
	}

	if adminOut == nil {
		return nil, admin.ErrAdminNotFound
	}

	out := admin.NewDeleteAdminOutput(adminOut, a.repo)

	if err := a.repo.Delete(id); err != nil {
		return nil, err
	}

	return out, nil

}
