package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg     interfaces.Logger
	listen string
	eChan  chan<- error

	server *http.Server
}

func New(
	lg interfaces.Logger,
	listen string,
	eChan chan<- error,
) *St {
	api := &St{
		lg:     lg,
		listen: listen,
		eChan:  eChan,
	}

	api.server = &http.Server{
		Addr:         listen,
		Handler:      promhttp.Handler(),
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 20 * time.Second,
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
