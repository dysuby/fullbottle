package dao

import (
	"bytes"
	"encoding/gob"
	"github.com/jinzhu/gorm"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
)

const (
	DEFAULT = 1 + iota
	ADMIN
)

type User struct {
	db.BasicModel
	Username  string `gorm:"type:varchar(24);not null" json:"username"`
	Password  string `gorm:"type:varchar(128);not null" json:"-"`
	Email     string `gorm:"type:varchar(128);not null" json:"email"`
	Role      int32  `gorm:"type:smallint;not null;default:1" json:"role"`
	AvatarFid string `gorm:"column:avatar_fid;type:varchar(64);not null" json:"avatarFid"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(*u)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (u *User) Unmarshal(b []byte) error {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(u)
	if err != nil {
		return err
	}

	return nil
}

func GetUsersById(id int64) (*User, error) {
	var user User
	if err := db.DB().Where("id = ?", id).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, common.NewDBError(err)
	}
	return &user, nil
}

func GetUsersByEmail(email string) (*User, error) {
	var user User
	if err := db.DB().Where("email = ?", email).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, common.NewDBError(err)
	}
	return &user, nil
}

func CreateUser(user *User) error {
	if err := db.DB().Create(user).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}

func UpdateUser(user *User, fields db.Fields) error {
	if err := db.DB().Model(user).Update(fields).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}
