package user_service

import "monitoring-system/server/src/internal/domain/user"

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) user.UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) GetByID(input *user.GetUserInput) (*user.User, error) {
	return u.repo.GetByID(input)
}

func (u *UserService) GetByEmail(email *user.GetUserByEmailInput) (*user.User, error) {
	return u.repo.GetByEmail(email)
}

func (u *UserService) Create(input *user.CreateUserInput) error {
	getUserByEmailInput := user.GetUserByEmailInput{Email: input.Email}
	if err := getUserByEmailInput.Validate(); err != nil {
		return err
	}

	userExists, err := u.repo.GetByEmail(&getUserByEmailInput)
	if err != nil {
		if err != user.ErrUserNotFound {
			return err
		}
	}

	if userExists != nil {
		return user.ErrUserAlreadyExists
	}

	return u.repo.Create(input)
}

func (u *UserService) RollbackCreate(input *user.CreateUserInput) error {
	return u.repo.Delete(&user.DeleteUserInput{ID: input.ID})
}

func (u *UserService) Update(input *user.UpdateUserInput) (backup *user.User, err error) {
	getUserInput := user.GetUserInput{ID: input.ID.String()}
	if err := getUserInput.Validate(); err != nil {
		return nil, err
	}

	out, err := u.repo.GetByID(&getUserInput)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, user.ErrUserNotFound
	}

	backup = out

	return backup, u.repo.Update(input)
}

func (u *UserService) RollbackUpdate(backup *user.User) error {
	rollbackInput := &user.UpdateUserInput{
		ID:    backup.ID,
		Name:  &backup.Name,
		Email: &backup.Email,
		Phone: backup.Phone,
	}
	return u.repo.Update(rollbackInput)
}

func (u *UserService) Delete(id *user.DeleteUserInput) error {
	getUserInput := user.GetUserInput{ID: id.ID.String()}
	if err := getUserInput.Validate(); err != nil {
		return err
	}

	out, err := u.repo.GetByID(&getUserInput)
	if err != nil {
		return err
	}

	if out == nil {
		return user.ErrUserNotFound
	}

	return u.repo.Delete(id)
}
