package common

import "github.com/micro/go-micro/v2/errors"

const (
	UserNotFound = 1000 + iota
	EmailExisted
	JwtError
	PasswordError

	DBConnError
	InternalError
)

func NewDBError(name string, err error) error {
	return errors.New(name, "Mysql Error: "+err.Error(), DBConnError)
}
