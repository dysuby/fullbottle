package config

import (
	"FullBottle/common/log"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/env"
)

type Config struct {
	Mysql struct {
		URL      string
		User     string
		Password string
		Database string
	}

	Redis struct {
		URL      string
		Password string
	}

	App struct {
		Secret string
	}

	Weed struct {
		Filer string
	}
}

var conf Config

func init() {
	c, err := config.NewConfig(config.WithSource(env.NewSource()))
	if err != nil {
		log.Fatalf(err, "Cannot load config")
	}
	if err = c.Scan(&conf); err != nil {
		log.Fatalf(err, "Config format error")
	}
}

func GetConfig() Config {
	return conf
}
