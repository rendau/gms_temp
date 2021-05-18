package core

import (
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg         interfaces.Logger
	cache      interfaces.Cache
	db         interfaces.Db
	sms        interfaces.Sms
	ws         interfaces.Ws
	noSmsCheck bool
	testing    bool

	Config *Config

	Dic     *Dic
	UsrType *UsrType

	Session *Session
	Usr     *Usr
}

func New(
	lg interfaces.Logger,
	cache interfaces.Cache,
	db interfaces.Db,
	sms interfaces.Sms,
	ws interfaces.Ws,
	noSmsCheck bool,
	testing bool,
) *St {
	c := &St{
		lg:         lg,
		cache:      cache,
		db:         db,
		sms:        sms,
		ws:         ws,
		noSmsCheck: noSmsCheck,
		testing:    testing,
	}

	c.Config = NewConfig(c)

	c.Dic = NewDic(c)
	c.UsrType = NewUsrType(c)

	c.Session = NewSession(c)
	c.Usr = NewUsr(c)

	return c
}
