package core

import (
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg    interfaces.Logger
	db    interfaces.Db
	cache interfaces.Cache

	Session *Session
}

func New(
	lg interfaces.Logger,
	db interfaces.Db,
	cache interfaces.Cache,
) *St {
	c := &St{
		lg:    lg,
		db:    db,
		cache: cache,
	}

	c.Session = NewSession(c)

	return c
}
