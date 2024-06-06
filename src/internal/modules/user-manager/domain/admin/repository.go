package admin

type AdminRepository interface {
	GetByID(input *GetAdminInput) (*Admin, error)
	GetByEmail(email *GetAdminByEmailInput) (*Admin, error)
	Create(input *CreateAdminInput) error
	Update(admin *UpdateAdminInput) error
	Delete(id *DeleteAdminInput) error
}
