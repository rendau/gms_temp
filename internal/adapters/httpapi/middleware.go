package httpapi

import (
	"context"
	"net/http"

	"github.com/rs/cors"
)

func (a *St) middleware(h http.Handler) http.Handler {
	h = cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           604800,
	}).Handler(h)

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
