package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
)

var db *gorm.DB

func init() {
	conf := config.C().Mysql

	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.URL, conf.Database)

	conn, err := gorm.Open("mysql", uri)
	if err != nil {
		log.WithError(err).Fatalf("Open db failed")
	}

	db = conn
}

func DB() *gorm.DB {
	return db
}
