package validator

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func ValidateStringLength(str string, min, max int) error {
	length := len(str)
	if length < min || length > max {
		return errors.New("string length is out of range")
	}
	return nil
}

func ValidateNumeric(str string) error {
	re := regexp.MustCompile(`^[0-9]+$`)
	if !re.MatchString(str) {
		return errors.New("must be a numeric value")
	}
	return nil
}
