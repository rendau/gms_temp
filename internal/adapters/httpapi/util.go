package httpapi

import (
	"encoding/json"
	"github.com/rendau/gms_temp/internal/domain/errs"
	"log"
	"net/http"
)

func (a *St) uRespondJSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Panicln("Fail to encode json obj", err)
	}
}

func (a *St) uHandleError(err error, w http.ResponseWriter) {
	if err != nil {
		switch cErr := err.(type) {
		case *errs.Err:
			a.uRespondJSON(w, ErrRepSt{
				ErrorCode: cErr.Error(),
			})
		default:
			a.uRespondJSON(w, ErrRepSt{
				ErrorCode: errs.ServiceNA.Error(),
			})
		}
	}
}
