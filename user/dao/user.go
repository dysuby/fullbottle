package dao

import (
	"github.com/vegchic/fullbottle/models"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	db = models.GetDB()
}

func GetUsersByQuery(query map[string]interface{}) ([]models.User, error) {
	var users []models.User
	if err := db.Where(query).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(user *models.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(user *models.User, fields map[string]interface{}) error {
	if err := db.Model(user).Update(fields).Error; err != nil {
		return err
	}
	return nil
}
