package core

import (
	"sync"

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

	wg     sync.WaitGroup
	stop   bool
	stopMu sync.RWMutex

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

func (c *St) StopAndWaitJobs() {
	c.stopMu.Lock()

	c.stop = true

	c.stopMu.Unlock()

	c.wg.Wait()
}

func (c *St) IsStopped() bool {
	c.stopMu.RLock()
	defer c.stopMu.RUnlock()

	return c.stop
}
