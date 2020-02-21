package dao

import (
	"github.com/vegchic/fullbottle/models"
)

func GetUsersByQuery(query map[string]interface{}) ([]models.User, error) {
	var users []models.User
	if err := models.DB().Where(query).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(user *models.User) error {
	if err := models.DB().Create(user).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(user *models.User, fields map[string]interface{}) error {
	if err := models.DB().Model(user).Update(fields).Error; err != nil {
		return err
	}
	return nil
}
