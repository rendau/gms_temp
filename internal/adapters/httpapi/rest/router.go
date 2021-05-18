package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *St) router() http.Handler {
	r := mux.NewRouter()

	// dic
	r.HandleFunc("/dic", a.hDicGet).Methods("GET")

	// config
	r.HandleFunc("/config", a.hConfigUpdate).Methods("PUT")

	// profile
	r.HandleFunc("/profile/send_validating_code", a.hProfileSendPhoneValidatingCode).Methods("POST")
	r.HandleFunc("/profile/auth", a.hProfileAuth).Methods("POST")
	r.HandleFunc("/profile/logout", a.hProfileLogout).Methods("POST")
	r.HandleFunc("/profile", a.hProfileGet).Methods("GET")
	r.HandleFunc("/profile", a.hProfileUpdate).Methods("PUT")
	r.HandleFunc("/profile/id", a.hProfileGetId).Methods("GET")

	// usrs
	r.HandleFunc("/usrs", a.hUsrList).Methods("GET")
	r.HandleFunc("/usrs", a.hUsrCreate).Methods("POST")
	r.HandleFunc("/usrs/{id:[0-9]+}", a.hUsrGet).Methods("GET")
	r.HandleFunc("/usrs/{id:[0-9]+}", a.hUsrUpdate).Methods("PUT")
	r.HandleFunc("/usrs/{id:[0-9]+}", a.hUsrDelete).Methods("DELETE")

	return a.middleware(r)
}
