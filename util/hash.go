package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Md5(str string) string {
	return BytesMd5([]byte(str))
}

func BytesMd5(bs []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bs))
}

func Sha256(str string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}

func Bcrypt(str string) string {
	password, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(password)
}

func BcryptCompare(expected string, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(expected), []byte(str))
	return err == nil
}
