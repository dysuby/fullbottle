package config

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/env"
	"github.com/vegchic/fullbottle/common/log"
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
		Upload struct {
			Secret string
		}
	}

	Server struct {
		Ip string
	}

	Weed struct {
		Master string
	}
}

var conf Config

func init() {
	c, err := config.NewConfig(config.WithSource(env.NewSource()))
	if err != nil {
		log.WithError(err).Fatalf("Cannot load config")
	}
	if err = c.Scan(&conf); err != nil {
		log.WithError(err).Fatalf("Config format error")
	}
}

func C() Config {
	return conf
}
