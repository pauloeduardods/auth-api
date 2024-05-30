package user

type user struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &user{
		repo: repo,
	}
}

func (u *user) GetByID(input *GetUserInput) (*User, error) {
	return u.repo.GetByID(input)
}

func (u *user) GetByEmail(email *GetUserByEmailInput) (*User, error) {
	return u.repo.GetByEmail(email)
}

func (u *user) Create(input *CreateUserInput) error {
	getUserByEmailInput, err := NewGetUserByEmailInput(input.Email)
	if err != nil {
		return err
	}

	userExists, err := u.repo.GetByEmail(&getUserByEmailInput)
	if err != nil {
		if err != ErrUserNotFound {
			return err
		}
	}

	if userExists != nil {
		return ErrUserAlreadyExists
	}

	return u.repo.Create(input)
}

func (u *user) RollbackCreate(input *CreateUserInput) error {
	return u.repo.Delete(&DeleteUserInput{ID: input.ID})
}

func (u *user) Update(input *UpdateUserInput) (backup *User, err error) {
	getUserInput, err := NewGetUserInput(input.ID.String())
	if err != nil {
		return nil, err
	}

	user, err := u.repo.GetByID(&getUserInput)
	if err != nil {
		return nil, err
	}

	backup = user

	if user == nil {
		return nil, ErrUserNotFound
	}

	return backup, u.repo.Update(input)
}

func (u *user) RollbackUpdate(backup *User) error {
	rollbackInput := &UpdateUserInput{
		ID:    backup.ID,
		Name:  &backup.Name,
		Email: &backup.Email,
		Phone: backup.Phone,
	}
	return u.repo.Update(rollbackInput)
}

func (u *user) Delete(id *DeleteUserInput) error {
	getUserInput, err := NewGetUserInput(id.ID.String())
	if err != nil {
		return err
	}

	user, err := u.repo.GetByID(&getUserInput)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	return u.repo.Delete(id)
}
