package config

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/env"
	"github.com/micro/go-micro/v2/util/log"
)

var conf config.Config

func init() {
	var err error
	conf, err = config.NewConfig(config.WithSource(env.NewSource()))
	if err != nil {
		log.Fatal(err)
	}
}

func GetSingleConfig(fields ...string) string {
	return conf.Get(fields...).String("")
}

func GetConfigMap(fields ...string) map[string]string {
	return conf.Get(fields...).StringMap(map[string]string{})
}
