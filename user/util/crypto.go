package util

import (
	"golang.org/x/crypto/bcrypt"
)

func PasswordCrypto(pwd string) string {
	password, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(password)
}

func ComparePassword(expected string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(expected), []byte(pwd))
	return err == nil
}
