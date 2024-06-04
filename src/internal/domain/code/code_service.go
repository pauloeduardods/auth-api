package code

type CodeService interface {
	GenerateAndSave(input GenerateAndSaveInput) (*Code, error)
	VerifyCode(input VerifyCodeInput) error
}
