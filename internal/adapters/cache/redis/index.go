package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/rendau/gms_temp/internal/interfaces"
	"time"
)

type St struct {
	lg interfaces.Logger
	r  *redis.Client
}

func NewRedisSt(lg interfaces.Logger, url, psw string, db int) *St {
	return &St{
		lg: lg,
		r: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: psw,
			DB:       db,
		}),
	}
}

func (c *St) Get(key string) ([]byte, bool, error) {
	data, err := c.r.Get(key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		c.lg.Errorw("Redis: fail to 'get'", err)
		return nil, false, err
	}
	return data, true, nil
}

func (c *St) Set(key string, value []byte, expiration time.Duration) error {
	err := c.r.Set(key, value, expiration).Err()
	if err != nil {
		c.lg.Errorw("Redis: fail to 'set'", err)
	}
	return err
}

func (c *St) Del(key string) error {
	err := c.r.Del(key).Err()
	if err != nil {
		c.lg.Errorw("Redis: fail to 'del'", err)
	}
	return err
}

func (c *St) Keys(pattern string) []string {
	var err error
	var cursor uint64
	var keys []string

	resKeys := make([]string, 0)
	for {
		keys, cursor, err = c.r.Scan(cursor, pattern, 30).Result()
		if err != nil {
			break
		}
		resKeys = append(resKeys, keys...)
		if cursor == 0 {
			break
		}
	}

	return resKeys
}
