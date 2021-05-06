package core

import (
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg    interfaces.Logger
	cache interfaces.Cache
	db    interfaces.Db
	sms   interfaces.Sms
	ws    interfaces.Ws

	Session *Session
}

func New(
	lg interfaces.Logger,
	cache interfaces.Cache,
	db interfaces.Db,
	sms interfaces.Sms,
	ws interfaces.Ws,
) *St {
	c := &St{
		lg:    lg,
		cache: cache,
		db:    db,
		sms:   sms,
		ws:    ws,
	}

	c.Session = NewSession(c)

	return c
}
