/*
Package rest GmsTemp API.

<br/><details>
	<summary>**Константы**</summary>
	```
	# Static file directories
	SFDUsrAva = "usr_avatar"

	# User types
	UsrTypeUndefined = 0
	UsrTypeAdmin     = 1

	# Notification types
	NfTypeRefreshProfile = "refresh-profile"
	NfTypeRefreshNumbers = "refresh-numbers"
	```
</details>

<details>
	<summary>**Работа с фото и файлами**</summary>
	[Документация файлового сервера](http://api.gms_temp.com/fs/doc/)<br/>
	Для заливки файлов на сервер необходимо указывать наименование папки (dir). Наименование нужно брать из констант
</details>

<details>
	<summary>**Websocket**</summary>

	websocket доступен по адресу `wss://ws.gms_temp.com?auth_token=<token>`
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
	r.HandleFunc("/mss/sms/balance/alarm", a.hSystemSmsBalanceAlarm).Methods("POST")
	r.HandleFunc("/mss/fs/filter_unused_files", a.hSystemFilterUnusedFiles).Methods("PUT")
	r.HandleFunc("/mss/cron/tick5m", a.hSystemCronTick5m).Methods("GET")
	r.HandleFunc("/mss/cron/tick15m", a.hSystemCronTick15m).Methods("GET")
	r.HandleFunc("/mss/cron/tick30m", a.hSystemCronTick30m).Methods("GET")

	// dic
	r.HandleFunc("/dic", a.hDicGet).Methods("GET")

	// config
	r.HandleFunc("/config", a.hConfigUpdate).Methods("PUT")

	// profile
	r.HandleFunc("/profile/send_validating_code", a.hProfileSendPhoneValidatingCode).Methods("POST")
	r.HandleFunc("/profile/auth", a.hProfileAuth).Methods("POST")
	r.HandleFunc("/profile/auth/token", a.hProfileAuthToken).Methods("POST")
	r.HandleFunc("/profile/reg", a.hProfileReg).Methods("POST")
	r.HandleFunc("/profile/logout", a.hProfileLogout).Methods("POST")
	r.HandleFunc("/profile", a.hProfileGet).Methods("GET")
	r.HandleFunc("/profile/numbers", a.hProfileGetNumbers).Methods("GET")
	r.HandleFunc("/profile", a.hProfileUpdate).Methods("PUT")
	r.HandleFunc("/profile/change_phone", a.hProfileChangePhone).Methods("PUT")
	r.HandleFunc("/profile/id", a.hProfileGetId).Methods("GET")

	// usr
	r.HandleFunc("/usr", a.hUsrList).Methods("GET")
	r.HandleFunc("/usr", a.hUsrCreate).Methods("POST")
	r.HandleFunc("/usr/{id:[0-9]+}", a.hUsrGet).Methods("GET")
	r.HandleFunc("/usr/{id:[0-9]+}", a.hUsrUpdate).Methods("PUT")
	r.HandleFunc("/usr/{id:[0-9]+}", a.hUsrDelete).Methods("DELETE")

	return a.middleware(r)
}
