package core

import (
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg    interfaces.Logger
	cache interfaces.Cache
	db    interfaces.Db
}

func New(lg interfaces.Logger, cache interfaces.Cache, db interfaces.Db) *St {
	c := &St{
		lg:    lg,
		cache: cache,
		db:    db,
	}

	return c
}
