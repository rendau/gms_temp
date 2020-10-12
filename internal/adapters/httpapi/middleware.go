package httpapi

import (
	"context"
	"github.com/rs/cors"
	"net/http"
)

func (a *St) middleware(h http.Handler) http.Handler {
	h = cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-Requested-With", "authorization"},
		MaxAge:         604800,
	}).Handler(h)
	h = a.mwRecovery(h)

	return h
}

func (a *St) mwRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancelCtx, cancel := context.WithCancel(r.Context())
		r = r.WithContext(cancelCtx)
		defer func() {
			if err := recover(); err != nil {
				cancel()
				w.WriteHeader(http.StatusInternalServerError)
				a.lg.Errorw(
					"Panic in http handler",
					err,
					"method", r.Method,
					"path", r.URL,
				)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
