package core

import (
	"context"
	"fmt"
	"time"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

const sessionContextKey = "user_session"

const sessionCacheKeyPattern = "user_session_%s"
const sessionCacheDuration = 20 * time.Minute

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

	cacheKey := fmt.Sprintf(sessionCacheKeyPattern, token)

	if cacheV := c.getFromCache(cacheKey); cacheV != nil {
		return cacheV
	}

	result.ID, result.TypeId, err = c.r.Usr.AuthByToken(ctx, token)
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

func (c *Session) SetToContext(ctx context.Context, ses *entities.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, ses)
}

func (c *Session) GetFromContext(ctx context.Context) *entities.Session {
	contextV := ctx.Value(sessionContextKey)
	if contextV == nil {
		return &entities.Session{}
	}

	switch ses := contextV.(type) {
	case *entities.Session:
		return ses
	default:
		c.r.lg.Fatal("wrong type of session in context")
		return &entities.Session{}
	}
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
	_ = c.r.cache.SetJsonObj(key, v, sessionCacheDuration)
}

func (c *Session) deleteUsrIdFromCache(id int64) bool {
	keys := c.r.cache.Keys(fmt.Sprintf(sessionCacheKeyPattern, "*"))

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
