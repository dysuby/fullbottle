package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/vegchic/fullbottle/config"
)

var client *redis.Client

func init() {
	conf := config.C().Redis

	client = redis.NewClient(&redis.Options{
		Addr:     conf.URL,
		Password: conf.Password,
		DB:       0,
	})
}

func Client() *redis.Client {
	return client
}
