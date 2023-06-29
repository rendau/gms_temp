package core

import (
	"context"
	"sync"

	"github.com/rendau/dop/adapters/cache"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/gms_temp/internal/adapters/repo"
)

type St struct {
	lg      logger.Lite
	cache   cache.Cache
	db      db.RDBContextTransaction
	repo    repo.Repo
	testing bool

	wg           sync.WaitGroup
	jobCtx       context.Context
	jobCtxCancel context.CancelFunc

	Config  *Config
	Dic     *Dic
	Session *Session
	System  *System
}

func New(
	lg logger.Lite,
	cache cache.Cache,
	db db.RDBContextTransaction,
	repo repo.Repo,
	testing bool,
) *St {
	jobCtx, jobCtxCancel := context.WithCancel(context.Background())

	c := &St{
		lg:           lg,
		cache:        cache,
		db:           db,
		repo:         repo,
		testing:      testing,
		jobCtx:       jobCtx,
		jobCtxCancel: jobCtxCancel,
	}

	c.Config = NewConfig(c)
	c.Dic = NewDic(c)
	c.Session = NewSession(c)
	c.System = NewSystem(c)

	return c
}

func (c *St) Start() {
}

func (c *St) IsStopped() bool {
	return c.jobCtx.Err() != nil
}

func (c *St) StopAndWaitJobs() {
	c.jobCtxCancel()
	c.wg.Wait()
}
