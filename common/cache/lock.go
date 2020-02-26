// Modify from https://github.com/bsm/redislock/blob/master/redislock.go

package cache

import (
	"errors"
	"github.com/vegchic/fullbottle/common"
	"math/rand"
	"strconv"
	"time"
	"unsafe"

	"github.com/go-redis/redis/v7"
)

var (
	luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL    = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
)

const DefaultTTL = 400 * time.Millisecond

type Lock struct {
	client *redis.Client
	key    string
	value  string
}

func (l *Lock) Key() string {
	return l.key
}

func (l *Lock) Token() string {
	return l.value
}

func (l *Lock) TTL() (time.Duration, error) {
	res, err := luaPTTL.Run(l.client, []string{l.key}, l.value).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if num := res.(int64); num > 0 {
		return time.Duration(num) * time.Millisecond, nil
	}
	return 0, nil
}

func (l *Lock) Refresh(ttl time.Duration) error {
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	res, err := luaRefresh.Run(l.client, []string{l.key}, l.value, ttlVal).Result()
	if err != nil {
		return err
	} else if res == 1 {
		return nil
	}

	return common.NewRedisError(errors.New("cannot refresh lock"))
}

func (l *Lock) Release() error {
	res, err := luaRelease.Run(l.client, []string{l.key}, l.value).Result()
	if err == redis.Nil {
		return common.NewRedisError(errors.New("cannot release lock"))
	} else if err != nil {
		return err
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return common.NewRedisError(errors.New("cannot release lock"))
	}
	return nil
}

func Obtain(key string, ttl time.Duration) (*Lock, error) {
	token := genToken(10)

	var timer *time.Timer
	for ddl := time.Now().Add(ttl); time.Now().Before(ddl); {
		ok, err := client.SetNX(key, token, ttl).Result()
		if err != nil {
			return nil, err
		} else if ok {
			return &Lock{client: client, key: key, value: token}, nil
		}

		if timer == nil {
			timer = time.NewTimer(time.Millisecond)
			defer timer.Stop()
		} else {
			timer.Reset(time.Millisecond)
		}

		select {
		case <-timer.C:
		}
	}

	return nil, common.NewRedisError(errors.New("cannot obtain lock"))
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func genToken(n int) string {
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
