package cache

import (
	"FullBottle/config"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

func init() {
	conf := config.GetConfig().Redis

	client = redis.NewClient(&redis.Options{
		Addr:     conf.URL,
		Password: conf.Password,
		DB:       0,
	})
}

func GetClient() *redis.Client {
	return client
}
