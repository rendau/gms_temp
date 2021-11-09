/*
Package rest GmsTemp API.

<br/><details>
	<summary>**Константы**</summary>
	```
	# Roles
	RoleGuest = "guest"
	RoleAdmin = "admin"
	```
</details>


    Schemes: https, http
    Host: api.gms_temp.com
    BasePath: /
    Version: 1.0.0

    Consumes:
    - application/json

    Produces:
    - application/json

    SecurityDefinitions:
      token:
         type: apiKey
         name: Authorization
         in: header
         description: "Пример: `Authorization: Bearer 2cf24dba5fb0a30e26e83b2`"

swagger:meta
*/
package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *St) router() http.Handler {
	r := mux.NewRouter()

	// doc
	r.HandleFunc("/doc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "doc/")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	r.PathPrefix("/doc/").Handler(http.StripPrefix("/doc/", http.FileServer(http.Dir("./doc/"))))

	// system
	r.HandleFunc("/mss/cron/tick5m", a.hSystemCronTick5m).Methods("GET")
	r.HandleFunc("/mss/cron/tick15m", a.hSystemCronTick15m).Methods("GET")
	r.HandleFunc("/mss/cron/tick30m", a.hSystemCronTick30m).Methods("GET")

	// dic
	r.HandleFunc("/dic", a.hDicGet).Methods("GET")

	// config
	r.HandleFunc("/config", a.hConfigUpdate).Methods("PUT")

	return a.middleware(r)
}
