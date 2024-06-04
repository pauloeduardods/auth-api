package email

type Email struct {
	To      string
	Subject string
	Body    string
}

func (e Email) Validate() error {
	if e.To == "" {
		return ErrEmailToEmpty
	}
	if e.Subject == "" {
		return ErrEmailSubjectEmpty
	}
	if e.Body == "" {
		return ErrEmailBodyEmpty
	}
	return nil
}
