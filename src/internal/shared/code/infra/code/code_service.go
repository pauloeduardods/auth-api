package code

import (
	"auth-api/src/internal/shared/code/domain/code"
	"auth-api/src/pkg/code_generator"
	"auth-api/src/pkg/logger"
	"context"
)

type CodeServiceImpl struct {
	codeRepo code.CodeRepository
	logger   logger.Logger
}

func NewCodeServiceImpl(repo code.CodeRepository, logger logger.Logger) code.CodeService {
	return &CodeServiceImpl{
		codeRepo: repo,
		logger:   logger,
	}
}

func (s *CodeServiceImpl) GenerateAndSave(ctx context.Context, input code.GenerateAndSaveInput) (*code.Code, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	codeValue, err := code_generator.GenerateCode(input.Length, input.CanContainLetters)
	if err != nil {
		return nil, err
	}

	code := &code.Code{
		Value:      codeValue,
		Identifier: input.Identifier,
		ExpiresAt:  input.ExpiresAt,
	}

	err = s.codeRepo.Save(ctx, code)
	if err != nil {
		return nil, err
	}

	return code, nil
}

func (s *CodeServiceImpl) VerifyCode(ctx context.Context, input code.VerifyCodeInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	codes, err := s.codeRepo.FindByIdentifier(ctx, input.Identifier)
	if err != nil {
		return err
	}

	if codes == nil || len(*codes) == 0 {
		return code.ErrCodeNotFound
	}

	for _, c := range *codes {
		if c.Value == input.Code {
			if c.IsExpired() {
				return code.ErrCodeExpired
			}
			if err := s.codeRepo.Delete(ctx, &c); err != nil {
				s.logger.Error("Error deleting code", err)
			}
			return nil
		}
	}

	return code.ErrInvalidCode
}
