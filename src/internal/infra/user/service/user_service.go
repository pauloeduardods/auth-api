package user_service

import "auth-api/src/internal/domain/user"

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

func (u *UserService) Create(input *user.CreateUserInput) (*user.CreateUserOutput, error) {
	getUserByEmailInput := user.GetUserByEmailInput{Email: input.Email}
	if err := getUserByEmailInput.Validate(); err != nil {
		return nil, err
	}

	userExists, err := u.repo.GetByEmail(&getUserByEmailInput)
	if err != nil {
		if err != user.ErrUserNotFound {
			return nil, err
		}
	}

	if userExists != nil {
		return nil, user.ErrUserAlreadyExists
	}

	out := &user.CreateUserOutput{
		ID: &input.ID,
	}

	if err := u.repo.Create(input); err != nil {
		return nil, err
	}

	return out, nil
}

func (u *UserService) Update(input *user.UpdateUserInput) (*user.UpdateUserOutput, error) {
	getUserInput := user.GetUserInput{ID: input.ID.String()}
	if err := getUserInput.Validate(); err != nil {
		return nil, err
	}

	userOut, err := u.repo.GetByID(&getUserInput)
	if err != nil {
		return nil, err
	}
	if userOut == nil {
		return nil, user.ErrUserNotFound
	}

	out := &user.UpdateUserOutput{
		Backup: userOut,
	}

	return out, u.repo.Update(input)
}

func (u *UserService) Delete(id *user.DeleteUserInput) (*user.DeleteUserOutput, error) {
	getUserInput := user.GetUserInput{ID: id.ID.String()}
	if err := getUserInput.Validate(); err != nil {
		return nil, err
	}

	userOut, err := u.repo.GetByID(&getUserInput)
	if err != nil {
		return nil, err
	}

	if userOut == nil {
		return nil, user.ErrUserNotFound
	}

	out := &user.DeleteUserOutput{
		Backup: userOut,
	}

	if err := u.repo.Delete(id); err != nil {
		return nil, err
	}

	return out, nil

}
