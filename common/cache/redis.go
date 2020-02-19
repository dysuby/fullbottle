package cache

import (
	"FullBottle/config"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

func init() {
	conf := config.GetConfigMap("redis")

	client = redis.NewClient(&redis.Options{
		Addr:     conf["url"],
		Password: conf["password"],
		DB:       0,
	})
}

func GetClient() *redis.Client {
	return client
}
