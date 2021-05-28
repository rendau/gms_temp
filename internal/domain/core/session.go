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

func (c *Session) Delete(id int64) {
	if c.deleteUsrIdFromCache(id) {
		c.r.Notification.SendRefreshProfile([]int64{id})
	}
}

func (c *Session) DeleteMany(ids []int64) {
	if len(ids) == 0 {
		return
	}

	var nfIds []int64

	for _, id := range ids {
		if c.deleteUsrIdFromCache(id) {
			nfIds = append(nfIds, id)
		}
	}

	c.r.Notification.SendRefreshProfile(nfIds)
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

func (c *Session) deleteUsrIdFromCache(id int64) bool {
	keys := c.r.cache.Keys(fmt.Sprintf(cacheKeyPattern, "*"))

	found := false

	for _, key := range keys {
		ses := entities.Session{}
		if sesFound, _ := c.r.cache.GetJsonObj(key, &ses); sesFound {
			if ses.ID == id {
				_ = c.r.cache.Del(key)

				found = true
			}
		}
	}

	return found
}
