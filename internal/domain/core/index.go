package core

import (
	"sync"

	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg      interfaces.Logger
	cache   interfaces.Cache
	db      interfaces.Db
	testing bool

	wg sync.WaitGroup

	Config *Config

	Dic *Dic

	System  *System
	Session *Session
}

func New(
	lg interfaces.Logger,
	cache interfaces.Cache,
	db interfaces.Db,
	testing bool,
) *St {
	c := &St{
		lg:      lg,
		cache:   cache,
		db:      db,
		testing: testing,
	}

	c.Config = NewConfig(c)

	c.Dic = NewDic(c)

	c.System = NewSystem(c)
	c.Session = NewSession(c)

	return c
}

func (c *St) Start() {
}

func (c *St) WaitJobs() {
	c.wg.Wait()
}
