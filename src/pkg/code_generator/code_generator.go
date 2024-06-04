package code_generator

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const digits = "0123456789"
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateCode(length int, canContainLetters bool) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be a positive integer")
	}

	var charset string
	if canContainLetters {
		charset = digits + letters
	} else {
		charset = digits
	}

	code := make([]byte, length)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}

	return string(code), nil
}
