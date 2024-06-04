package code_service

import (
	"auth-api/src/internal/domain/code"
	"auth-api/src/pkg/code_generator"
	"time"
)

type CodeServiceImpl struct {
	CodeRepo code.CodeRepository
}

func NewCodeServiceImpl(repo code.CodeRepository) code.CodeService {
	return &CodeServiceImpl{
		CodeRepo: repo,
	}
}

func (s *CodeServiceImpl) GenerateAndSave(identifier string, expiresAt time.Time, length int, canContainLetters bool) (*code.Code, error) {
	codeValue, err := code_generator.GenerateCode(length, canContainLetters)
	if err != nil {
		return nil, err
	}

	code := &code.Code{
		Value:      codeValue,
		ExpiresAt:  expiresAt,
		Identifier: identifier,
	}

	err = s.CodeRepo.Save(code)
	if err != nil {
		return nil, err
	}

	return code, nil
}

func (s *CodeServiceImpl) VerifyCode(identifier, value string) error {
	codes, err := s.CodeRepo.FindByIdentifier(identifier)
	if err != nil {
		return err
	}

	if codes == nil || len(*codes) == 0 {
		return code.ErrCodeNotFound
	}

	for _, c := range *codes {
		if !c.IsExpired() && c.Value == value {
			return nil
		}
	}

	return code.ErrInvalidCode
}
