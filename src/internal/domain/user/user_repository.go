package user

type UserRepository interface {
	GetByID(input *GetUserInput) (*User, error)
	GetByEmail(email *GetUserByEmailInput) (*User, error)
	Create(input *CreateUserInput) error
	Update(user *UpdateUserInput) error
	Delete(id *DeleteUserInput) error
}
