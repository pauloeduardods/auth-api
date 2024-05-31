package user

type UserService interface {
	GetByID(input *GetUserInput) (*User, error)
	GetByEmail(input *GetUserByEmailInput) (*User, error)
	Create(input *CreateUserInput) (*CreateUserOutput, error)
	Update(input *UpdateUserInput) (*UpdateUserOutput, error)
	Delete(input *DeleteUserInput) (*DeleteUserOutput, error)
}
