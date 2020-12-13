package usecases

import (
	"context"
	"encoding/json"
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
	raw, ok, err := u.cc.Get(key)
	if err != nil {
		return nil
	}

	if !ok {
		return nil
	}

	result := &entities.Session{}

	err = json.Unmarshal(raw, result)
	if err != nil {
		return nil
	}

	return result
}

func (u *St) sessionSetToCache(key string, v *entities.Session) {
	raw, err := json.Marshal(v)
	if err != nil {
		return
	}

	err = u.cc.Set(key, raw, cacheDuration)
	if err != nil {
		return
	}
}

func (u *St) sessionDeleteUsrIdFromCache(id int64) {
	var err error
	var ses entities.Session
	var raw []byte
	var found bool

	keys := u.cc.Keys(fmt.Sprintf(cacheKeyPattern, "*"))

	for _, key := range keys {
		raw, found, _ = u.cc.Get(key)
		if found {
			err = json.Unmarshal(raw, &ses)
			if err == nil {
				if ses.ID == id {
					_ = u.cc.Del(key)
				}
			}
		}
	}
}
