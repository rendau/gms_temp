package core

import (
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg    interfaces.Logger
	cache interfaces.Cache
	db    interfaces.Db

	Session *Session
}

func New(
	lg interfaces.Logger,
	cache interfaces.Cache,
	db interfaces.Db,
) *St {
	c := &St{
		lg:    lg,
		cache: cache,
		db:    db,
	}

	c.Session = NewSession(c)

	return c
}
