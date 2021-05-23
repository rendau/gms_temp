package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

// swagger:route GET /dic dic hDicGet
// Все справочники.
// Это АПИ нужно вызывать при каждом запуске приложения (фронта)
// Responses:
//   200: dicRep
//   400: errRep
func (a *St) hDicGet(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hDicGet
	type docReqSt struct {
		// Хэш сумма значении ответа. Используется для сверки версии ответа от сервера.
		// В первый раз нужно передать пустую строку, потом нужно передать то что вернул сервер в ответе.
		// in:query
		Hs string `json:"hs"`
	}

	// `hs` - фронту нужно сохранить значение этого поля и передать в параметре следующего запроса
	// <br/>`data` - может быть null, если __hs__ переданный в запросе совпадает с серверным __hs__, в этом случае фронт должен использовать предыдущее значение ответа от сервера
	// swagger:response dicRep
	type docRepSt struct {
		// in:body
		Body struct {
			Hs   string             `json:"hs"`
			Data entities.DicDataSt `json:"data"`
		}
	}

	reqQuery := r.URL.Query()

	rhs := reqQuery.Get("hs")

	hs, data, err := a.ucs.DicGetJson(a.uGetRequestContext(r), rhs)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, struct {
		Hs   string          `json:"hs"`
		Data json.RawMessage `json:"data"`
	}{
		Hs:   hs,
		Data: data,
	})
}
