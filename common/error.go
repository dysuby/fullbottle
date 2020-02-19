package common

const (
	Success = 1000 + iota

	UserNotFound
	ArgumentError
	EmailExisted
	JwtError
	PasswordError

	DBConnError
	InternalError
)
