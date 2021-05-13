package core

import (
	"context"
	"fmt"
	"time"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

const cacheKeyPattern = "user_session_%s"
const cacheDuration = 20 * time.Minute

type Session struct {
	r *St
}

func NewSession(r *St) *Session {
	return &Session{r: r}
}

func (c *Session) Get(ctx context.Context, token string) *entities.Session {
	var err error

	result := &entities.Session{}

	if token == "" {
		return result
	}

	cacheKey := fmt.Sprintf(cacheKeyPattern, token)

	if cacheV := c.getFromCache(cacheKey); cacheV != nil {
		return cacheV
	}

	result.ID, err = c.r.Usr.AuthByToken(ctx, token)
	if err != nil {
		return result
	}

	result.TypeId, err = c.r.Usr.GetTypeId(ctx, result.ID)
	if err != nil {
		return result
	}

	c.setToCache(cacheKey, result)

	return result
}

func (c *Session) Delete(ctx context.Context, id int64) {
	c.deleteUsrIdFromCache(id)
}

func (c *Session) getFromCache(key string) *entities.Session {
	result := &entities.Session{}

	ok, err := c.r.cache.GetJsonObj(key, &result)
	if err != nil {
		return nil
	}
	if !ok {
		return nil
	}

	return result
}

func (c *Session) setToCache(key string, v *entities.Session) {
	_ = c.r.cache.SetJsonObj(key, v, cacheDuration)
}

func (c *Session) deleteUsrIdFromCache(id int64) {
	keys := c.r.cache.Keys(fmt.Sprintf(cacheKeyPattern, "*"))

	for _, key := range keys {
		ses := entities.Session{}
		found, _ := c.r.cache.GetJsonObj(key, &ses)
		if found {
			if ses.ID == id {
				_ = c.r.cache.Del(key)
			}
		}
	}
}
