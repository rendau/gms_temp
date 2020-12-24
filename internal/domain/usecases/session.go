package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/rendau/gms_temp/internal/domain/errs"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

const cacheKeyPattern = "user_session_%s"
const cacheDuration = 20 * time.Minute

func (u *St) SessionGet(ctx context.Context, token string) *entities.Session {
	result := &entities.Session{}

	if token == "" {
		return result
	}

	cacheKey := fmt.Sprintf(cacheKeyPattern, token)

	if cacheV := u.sessionGetFromCache(cacheKey); cacheV != nil {
		return cacheV
	}

	// try to auth in db with the token, then set to cache

	return result
}

func (u *St) SesRequireAuth(ses *entities.Session) error {
	if ses.ID == 0 {
		return errs.NotAuthorized
	}
	return nil
}

func (u *St) SessionDelete(ctx context.Context, id int64) {
	u.sessionDeleteUsrIdFromCache(id)
}

func (u *St) sessionGetFromCache(key string) *entities.Session {
	result := &entities.Session{}

	ok, err := u.cache.GetJsonObj(key, &result)
	if err != nil {
		return nil
	}
	if !ok {
		return nil
	}

	return result
}

func (u *St) sessionSetToCache(key string, v *entities.Session) {
	_ = u.cache.SetJsonObj(key, v, cacheDuration)
}

func (u *St) sessionDeleteUsrIdFromCache(id int64) {
	keys := u.cache.Keys(fmt.Sprintf(cacheKeyPattern, "*"))

	for _, key := range keys {
		ses := entities.Session{}
		found, _ := u.cache.GetJsonObj(key, &ses)
		if found {
			if ses.ID == id {
				_ = u.cache.Del(key)
			}
		}
	}
}
