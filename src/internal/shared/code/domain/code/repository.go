package code

import "context"

type CodeRepository interface {
	Save(ctx context.Context, code *Code) error
	FindCode(ctx context.Context, identifier, code string) (*Code, error)
	Delete(ctx context.Context, code *Code) error
}
