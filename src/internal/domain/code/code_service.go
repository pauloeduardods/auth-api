package code

import "time"

type CodeService interface {
	GenerateAndSave(identifier string, expiresAt time.Time, length int, canContainLetters bool) (*Code, error)
	VerifyCode(identifier, code string) error
}
