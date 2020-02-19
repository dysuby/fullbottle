package models

import (
	"FullBottle/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2/util/log"
)

var db *gorm.DB

func init() {
	conf := config.GetConfigMap("mysql")
	user := conf["user"]
	url := conf["url"]
	password := conf["password"]
	database := conf["database"]

	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, url, database)
	conn, err := gorm.Open("mysql", uri)
	if err != nil {
		log.Fatalf("Mysql uri: %s, err: %v", uri, err)
	}
	db = conn
}

func GetDB() *gorm.DB {
	return db
}
