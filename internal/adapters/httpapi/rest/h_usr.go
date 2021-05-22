package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/util"
)

func (a *St) hUsrList(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /usrs usrs usr_list
	// Список пользователей.
	// Список пользователей
	// <br/>Доступен только для: `UsrTypeAdmin`
	//   Security:
	//     token:
	//   Responses:
	//     200: usrListRep
	//     400: error_reply

	// swagger:parameters usr_list
	type usrListReq struct {
		// in:query
		entities.PaginationParams

		// in:query
		TypeId int `json:"type_id"`

		// in:query
		OnlyCount int `json:"only_count"`

		// in:query
		Search int `json:"search"`
	}

	// swagger:response usrListRep
	type usrListRep struct {
		// in:body
		Body struct {
			DocPaginatedListRepSt
			Results []*entities.UsrListSt `json:"results"`
		}
	}

	qPars := r.URL.Query()

	pars := &entities.UsrListParsSt{
		TypeId:    a.uQpParseInt(qPars, "type_id"),
		OnlyCount: a.uQpParseBoolV(qPars, "only_count"),
		Search:    a.uQpParseString(qPars, "search"),
	}

	if idPar := a.uQpParseInt64(qPars, "id"); idPar != nil {
		pars.Ids = util.NewSliceInt64(*idPar)
	}

	offset, limit, page := a.uExtractPaginationPars(qPars)
	pars.Offset = offset
	pars.Limit = limit

	paginated := pars.Limit > 0

	result, tCount, err := a.ucs.UsrList(a.uGetRequestContext(r), pars)
	if a.uHandleError(err, r, w) {
		return
	}

	if paginated {
		a.uRespondJSON(w, 0, &PaginatedListRepSt{
			Page:       page,
			PageSize:   limit,
			TotalCount: tCount,
			Results:    result,
		})
	} else {
		a.uRespondJSON(w, 0, result)
	}
}

func (a *St) hUsrCreate(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /usrs usrs usr_create
	// Создать пользователя.
	// Создать пользователя
	// <br/>Доступен только для: `UsrTypeAdmin`
	//   Security:
	//     token:
	//   Responses:
	//     200: usrCreateRep
	//     400: error_reply

	// swagger:parameters usr_create
	type usrCreateReq struct {
		// in: body
		Body entities.UsrCUSt
	}

	// `id` - id созданной записи
	// swagger:response usrCreateRep
	type usrCreateRep struct {
		// in:body
		Body struct {
			Id int64 `json:"id"`
		}
	}

	reqObj := &entities.UsrCUSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	newId, err := a.ucs.UsrCreate(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, map[string]int64{"id": newId})
}

func (a *St) hUsrGet(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /usrs/{id} usrs usr_get
	// Пользователь.
	// Пользователь
	//   Security:
	//     token:
	//   Responses:
	//     200: usrGetRep
	//     400: error_reply

	// swagger:parameters usr_get
	type usrGetReq struct {
		// in: path
		Id int64 `json:"id"`
	}

	// swagger:response usrGetRep
	type usrGetRep struct {
		// in:body
		Body entities.UsrSt
	}

	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	// qPars := r.URL.Query()

	result, err := a.ucs.UsrGet(a.uGetRequestContext(r), id)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, 0, result)
}

func (a *St) hUsrUpdate(w http.ResponseWriter, r *http.Request) {
	// swagger:route PUT /usrs/{id} usrs usr_update
	// Изменить пользователя.
	// Изменить пользователя
	// <br/>Доступен только для: `UsrTypeAdmin`
	//   Security:
	//     token:
	//   Responses:
	//     200:
	//     400: error_reply

	// swagger:parameters usr_update
	type usrUpdateReq struct {
		// in: path
		Id int64 `json:"id"`

		// in: body
		Body entities.UsrCUSt
	}

	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	reqObj := &entities.UsrCUSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.UsrUpdate(a.uGetRequestContext(r), id, reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

func (a *St) hUsrDelete(w http.ResponseWriter, r *http.Request) {
	// swagger:route DELETE /usrs/{id} usrs usr_delete
	// Удалить пользователя.
	// Удалить пользователя
	// <br/>Доступен только для: `UsrTypeAdmin`
	//   Security:
	//     token:
	//   Responses:
	//     200:
	//     400: error_reply

	// swagger:parameters usr_delete
	type usrDeleteReq struct {
		// in: path
		Id int64 `json:"id"`
	}

	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	err := a.ucs.UsrDelete(a.uGetRequestContext(r), id)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
