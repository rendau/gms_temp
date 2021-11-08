package rest

import (
	"net/http"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

// swagger:route POST /profile/send_validating_code profile hProfileSendPhoneValidatingCode
// Отправить СМС код на номер.
// Responses:
//   200:
//   400: errRep
func (a *St) hProfileSendPhoneValidatingCode(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hProfileSendPhoneValidatingCode
	type docReqSt struct {
		// `err_ne` если передать **true** - то вернет в ответе ошибку если номера нет в базе
		// in: body
		Body struct {
			Phone string `json:"phone"`
			ErrNE bool   `json:"err_ne"`
		}
	}

	reqObj := &(struct {
		Phone string `json:"phone"`
		ErrNE bool   `json:"err_ne"`
	}{})
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.ProfileSendPhoneValidatingCode(a.uGetRequestContext(r), reqObj.Phone, reqObj.ErrNE)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

// swagger:route POST /profile/auth profile hProfileAuth
// Авторизация.
// Responses:
//   200: profileAuthRep
//   400: errRep
func (a *St) hProfileAuth(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hProfileAuth
	type docReqSt struct {
		// in: body
		Body entities.PhoneAndSmsCodeSt
	}

	// swagger:response profileAuthRep
	type docRepSt struct {
		// in:body
		Body struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
	}

	reqObj := &entities.PhoneAndSmsCodeSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	accessToken, refreshToken, err := a.ucs.ProfileAuth(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// swagger:route POST /profile/auth/token profile hProfileAuthToken
// Авторизация.
// Responses:
//   200: profileAuthTokenRep
//   400: errRep
func (a *St) hProfileAuthToken(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hProfileAuthToken
	type docReqSt struct {
		// in: body
		Body entities.PhoneAndSmsCodeSt
	}

	// swagger:response profileAuthTokenRep
	type docRepSt struct {
		// in:body
		Body struct {
			AccessToken string `json:"access_token"`
		}
	}

	reqObj := map[string]string{}
	if !a.uParseRequestJSON(w, r, &reqObj) {
		return
	}

	accessToken, err := a.ucs.ProfileAuthByRefreshToken(a.uGetRequestContext(r), reqObj["refresh_token"])
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, map[string]string{
		"access_token": accessToken,
	})
}

// swagger:route POST /profile/reg profile hProfileReg
// Регистрация.
// Responses:
//   200: profileRegRep
//   400: errRep
func (a *St) hProfileReg(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hProfileReg
	type docReqSt struct {
		// `type_id` игнорируется
		// in: body
		Body entities.UsrRegReqSt
	}

	// swagger:response profileRegRep
	type docRepSt struct {
		// in:body
		Body struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
	}

	reqObj := &entities.UsrRegReqSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	accessToken, refreshToken, err := a.ucs.ProfileReg(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// swagger:route POST /profile/logout profile hProfileLogout
// Разлогиниться.
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
func (a *St) hProfileLogout(w http.ResponseWriter, r *http.Request) {
	err := a.ucs.ProfileLogout(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

// swagger:route GET /profile profile hProfileGet
// Объект профиля.
// Security:
//   token:
// Responses:
//   200: profileRep
//   400: errRep
func (a *St) hProfileGet(w http.ResponseWriter, r *http.Request) {
	// swagger:response profileRep
	type docRepSt struct {
		// in:body
		Body entities.UsrProfileSt
	}

	profile, err := a.ucs.ProfileGet(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, profile)
}

// swagger:route GET /profile/numbers profile hProfileGetNumbers
// Новые цифры профиля (badge).
// Используется чтоб показать на фронте сколько уведомлении клиент пропустил
// Security:
//   token:
// Responses:
//   200: profileNumbersRep
//   400: errRep
func (a *St) hProfileGetNumbers(w http.ResponseWriter, r *http.Request) {
	// swagger:response profileNumbersRep
	type docRepSt struct {
		// in:body
		Body entities.UsrNumbersSt
	}

	result, err := a.ucs.ProfileGetNumbers(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, result)
}

// swagger:route PUT /profile profile hProfileUpdate
// Изменить данные профиля.
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
func (a *St) hProfileUpdate(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hProfileUpdate
	type docReqSt struct {
		// `type_id`, `phone` игнорируются
		// in: body
		Body entities.UsrCUSt
	}

	reqObj := &entities.UsrCUSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.ProfileUpdate(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

// swagger:route PUT /profile/change_phone profile hProfileChangePhone
// Изменить основной номер.
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
func (a *St) hProfileChangePhone(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hProfileChangePhone
	type docReqSt struct {
		// `phone` - это новый номер
		// <br/>`sms_code` должен быть отправлен на новый номер
		// in: body
		Body entities.PhoneAndSmsCodeSt
	}

	reqObj := &entities.PhoneAndSmsCodeSt{}
	if !a.uParseRequestJSON(w, r, &reqObj) {
		return
	}

	err := a.ucs.ProfileChangePhone(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}

// swagger:route GET /profile/id profile hProfileGetId
// Id профиля.
// В основном используется для интеграции с другими сервисами в бэкенде
// Security:
//   token:
// Responses:
//   200: profileIdRep
//   400: errRep
func (a *St) hProfileGetId(w http.ResponseWriter, r *http.Request) {
	// swagger:response profileIdRep
	type docRepSt struct {
		// in:body
		Body struct {
			Id int64 `json:"id"`
		}
	}

	usrId, err := a.ucs.ProfileGetId(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, map[string]int64{
		"id": usrId,
	})
}

func (a *St) hProfileDelete(w http.ResponseWriter, r *http.Request) {
	err := a.ucs.ProfileDelete(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
