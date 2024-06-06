package code

import "context"

type CodeService interface {
	GenerateAndSave(ctx context.Context, input GenerateAndSaveInput) (*Code, error)
	VerifyCode(ctx context.Context, input VerifyCodeInput) error
}
