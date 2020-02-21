package dao

import (
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/common/log"
)

const (
	DEFAULT = 1 + iota
	ADMIN
)

type User struct {
	db.BasicModel
	Username  string `gorm:"type:varchar(24);not null"`
	Password  string `gorm:"type:varchar(128);not null"`
	Email     string `gorm:"type:varchar(128);not null"`
	Role      int32  `gorm:"type:smallint;not null;default:1"`
	AvatarFid string `gorm:"type:varchar(64);not null"`
}

func GetUsersByQuery(query map[string]interface{}) ([]User, error) {
	var users []User
	if err := db.DB().Where(query).Find(&users).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return users, nil
}

func CreateUser(user *User) error {
	if err := db.DB().Create(user).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}

func UpdateUser(user *User, fields map[string]interface{}) error {
	if err := db.DB().Model(user).Update(fields).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}
