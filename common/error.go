package common

import (
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/config"
)

const (
	UserNotFound = 1000 + iota
	EmailExisted
	JwtError
	PasswordError
	EmptyAvatarError

	FileFetchError
	FileUploadError

	DBConnError
	InternalError
	WeedError
)

func NewDBError(err error) error {
	return errors.New(config.DBName, "Mysql Error: "+err.Error(), DBConnError)
}

func NewWeedError(err error) error {
	return errors.New(config.WeedName, "Weed Error: "+err.Error(), WeedError)
}
