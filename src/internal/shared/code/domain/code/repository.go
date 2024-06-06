package code

import "context"

type CodeRepository interface {
	Save(ctx context.Context, code *Code) error
	FindByIdentifier(ctx context.Context, identifier string) (*[]Code, error)
	Delete(ctx context.Context, code *Code) error
}
