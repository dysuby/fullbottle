package models

import (
	"FullBottle/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"FullBottle/common/log"
)

var db *gorm.DB

func init() {
	conf := config.GetConfig().Mysql

	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.URL, conf.Database)

	conn, err := gorm.Open("mysql", uri)
	if err != nil {
		log.Fatalf(err, "Open db failed")
	}

	db = conn
}

func GetDB() *gorm.DB {
	return db
}
