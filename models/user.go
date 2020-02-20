package models

const (
	DEFAULT = 1 + iota
	ADMIN
)

type User struct {
	BasicModel
	Username  string `gorm:"type:varchar(24);not null"`
	Password  string `gorm:"type:varchar(128);not null"`
	Email     string `gorm:"type:varchar(128);not null"`
	Role      int32  `gorm:"type:smallint;not null;default:1"`
	AvatarUri string `gorm:"type:varchar(64);not null"`
}
