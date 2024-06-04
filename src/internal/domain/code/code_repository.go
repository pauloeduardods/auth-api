package code

type CodeRepository interface {
	Save(code *Code) error
	FindByIdentifier(identifier string) (*[]Code, error)
	Delete(code *Code) error
}
