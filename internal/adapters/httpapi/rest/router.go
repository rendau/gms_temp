/*
Package rest GmsTemp API.

<br/><details>
	<summary>**Константы**</summary>
	```
	AppUrl = "http://gms_temp.com"

	UsrTypeUndefined = 0
	UsrTypeAdmin = 1
	```
</details>

<details>
	<summary>**Работа с фото и файлами**</summary>
	[Документация файлового сервера](http://gms_temp.com/api/fs/doc/)<br/>
</details>

<details>
	<summary>**Websocket**</summary>

	websocket доступен по адресу `wss://gms_temp.com/api/websocket?auth_token=<token>`. По вебсокету приходят такие же данные как и в пушах.
</details>


    Schemes: https, http
    Host: gms_temp.com
    BasePath: /api
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
         description: "Пример: `Authorization: 2cf24dba5fb0a30e26e83b2`"

swagger:meta
*/
package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	mmStd "github.com/slok/go-http-metrics/middleware/std"
)

func (a *St) router() http.Handler {
	r := mux.NewRouter()

	mh := func(h http.HandlerFunc, id string) http.Handler {
		return h
	}

	if a.withMetrics {
		mm := middleware.New(middleware.Config{
			Recorder: prometheus.NewRecorder(prometheus.Config{}),
		})

		mh = func(h http.HandlerFunc, id string) http.Handler {
			return mmStd.Handler(id, mm, h)
		}
	}

	// doc
	r.HandleFunc("/doc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "doc/")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	r.PathPrefix("/doc/").Handler(http.StripPrefix("/doc/", http.FileServer(http.Dir("./doc/"))))

	// system
	r.HandleFunc("/mss/sms/balance/alarm", a.hSystemSmsBalanceAlarm).Methods("POST")
	r.HandleFunc("/mss/fs/filter_unused_files", a.hSystemFilterUnusedFiles).Methods("PUT")
	r.HandleFunc("/mss/cron/tick5m", a.hSystemCronTick5m).Methods("GET")
	r.HandleFunc("/mss/cron/tick15m", a.hSystemCronTick15m).Methods("GET")
	r.HandleFunc("/mss/cron/tick30m", a.hSystemCronTick30m).Methods("GET")

	// dic
	r.Handle("/dic", mh(a.hDicGet, "/dic")).Methods("GET")

	// config
	r.HandleFunc("/config", a.hConfigUpdate).Methods("PUT")

	// profile
	r.Handle("/profile/send_validating_code", mh(a.hProfileSendPhoneValidatingCode, "/profile/send_validating_code")).Methods("POST")
	r.Handle("/profile/auth", mh(a.hProfileAuth, "/profile/auth")).Methods("POST")
	r.Handle("/profile/reg", mh(a.hProfileReg, "/profile/reg")).Methods("POST")
	r.Handle("/profile/logout", mh(a.hProfileLogout, "/profile/logout")).Methods("POST")
	r.Handle("/profile", mh(a.hProfileGet, "/profile")).Methods("GET")
	r.Handle("/profile/numbers", mh(a.hProfileGetNumbers, "/profile/numbers")).Methods("GET")
	r.Handle("/profile", mh(a.hProfileUpdate, "/profile")).Methods("PUT")
	r.Handle("/profile/change_phone", mh(a.hProfileChangePhone, "/profile/change_phone")).Methods("PUT")
	r.Handle("/profile/id", mh(a.hProfileGetId, "/profile/id")).Methods("GET")

	// usrs
	r.Handle("/usrs", mh(a.hUsrList, "/usrs")).Methods("GET")
	r.Handle("/usrs", mh(a.hUsrCreate, "/usrs")).Methods("POST")
	r.Handle("/usrs/{id:[0-9]+}", mh(a.hUsrGet, "/usrs/:id")).Methods("GET")
	r.Handle("/usrs/{id:[0-9]+}", mh(a.hUsrUpdate, "/usrs/:id")).Methods("PUT")
	r.Handle("/usrs/{id:[0-9]+}", mh(a.hUsrDelete, "/usrs/:id")).Methods("DELETE")

	return a.middleware(r)
}
