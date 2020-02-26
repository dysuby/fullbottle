package common

import (
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/config"
)

const (
	NotFoundError = 1000 + iota
	ExistedError
	ConflictError
	JwtError
	PasswordError
	EmptyAvatarError
	ExceedError

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

func NewRedisError(err error) error {
	return errors.New(config.RedisName, "Redis Error: "+err.Error(), WeedError)

}
