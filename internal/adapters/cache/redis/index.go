package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg interfaces.Logger
	r  *redis.Client
}

func New(lg interfaces.Logger, url, psw string, db int) *St {
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

func (c *St) GetJsonObj(key string, dst interface{}) (bool, error) {
	dataRaw, ok, err := c.Get(key)
	if err != nil || !ok {
		return ok, err
	}

	err = json.Unmarshal(dataRaw, dst)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *St) Set(key string, value []byte, expiration time.Duration) error {
	err := c.r.Set(key, value, expiration).Err()
	if err != nil {
		c.lg.Errorw("Redis: fail to 'set'", err)
	}

	return err
}

func (c *St) SetJsonObj(key string, value interface{}, expiration time.Duration) error {
	dataRaw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Set(key, dataRaw, expiration)
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
