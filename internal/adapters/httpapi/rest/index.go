package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/rendau/gms_temp/internal/domain/usecases"
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	cors        bool
	lg          interfaces.Logger
	ucs         *usecases.St
	withMetrics bool
	eChan       chan<- error

	server *http.Server
}

func New(
	cors bool,
	lg interfaces.Logger,
	listen string,
	ucs *usecases.St,
	withMetrics bool,
	eChan chan<- error,
) *St {
	api := &St{
		cors:        cors,
		lg:          lg,
		ucs:         ucs,
		withMetrics: withMetrics,
		eChan:       eChan,
	}

	api.server = &http.Server{
		Addr:              listen,
		Handler:           api.router(),
		ReadTimeout:       2 * time.Minute,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return api
}

func (a *St) Start() {
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.lg.Errorw("Http server closed", err)
			a.eChan <- err
		}
	}()
}

func (a *St) Shutdown(timeout time.Duration) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()

	return a.server.Shutdown(ctx)
}
