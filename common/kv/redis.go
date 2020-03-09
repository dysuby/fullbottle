package kv

import (
	"github.com/go-redis/redis/v7"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
	"time"
)

var client *redis.Client

var luaRefreshValue = redis.NewScript("local ttl = redis.call('ttl', KEYS[1]) if ttl > 0 then return redis.call('SETEX', KEYS[1], ttl, ARGV[1]) else return 0 end")

func init() {
	conf := config.C().Redis

	client = redis.NewClient(&redis.Options{
		Addr:     conf.URL,
		Password: conf.Password,
		DB:       0,
	})
}

func C() *redis.Client {
	return client
}

type Marshaller interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

func GetM(key string, m Marshaller) error {
	if raw, err := client.Get(key).Bytes(); err == nil {
		if err = m.Unmarshal(raw); err != nil {
			log.WithError(err).Infof("Unmarshal failed")
			client.Del(key)
		}
		return nil
	} else if err == redis.Nil {
		return common.NewRedisError(err)
	} else {
		log.WithError(err).Infof("Redis failed")
		return common.NewRedisError(err)
	}
}

func SetM(key string, m Marshaller, exp time.Duration) error {
	if b, err := m.Marshal(); err != nil {
		log.WithError(err).Infof("Marshal failed")
		return common.NewRedisError(err)
	} else if err = client.Do("SETEX", key, exp.Milliseconds(), b).Err(); err != nil {
		log.WithError(err).Infof("Redis failed")
		return common.NewRedisError(err)
	}
	return nil
}

func RefreshMValue(key string, m Marshaller) error {
	if b, err := m.Marshal(); err != nil {
		log.WithError(err).Infof("Marshal failed")
		return common.NewRedisError(err)
	} else if err = luaRefreshValue.Run(client, []string{key}, b).Err(); err != nil {
		log.WithError(err).Infof("Redis failed")
		return common.NewRedisError(err)
	}
	return nil
}

func Del(key string) error {
	if err := client.Del(key).Err(); err != nil {
		log.WithError(err).Infof("Redis failed")
		return common.NewRedisError(err)
	}
	return nil
}
