package swagger

import (
	"github.com/rendau/gms_temp/internal/domain/entities"
)

/*
swagger:route GET /dic dic get_dic
Все справочники.
Responses:
  200: dicRep
*/

// swagger:parameters get_dic
type dicReqQuery struct {
	// Хэш сумма значении ответа. Используется для сверки версии ответа от сервера.
	// В первый раз нужно передать пустую строку, потом нужно передать то что вернул сервер в ответе.
	// in:query
	Hs string `json:"hs"`
}

/*
`data` - может быть null, если __hs__ переданный в запросе совпадает с серверным __hs__, в этом случае фронт должен использовать предыдущее значение ответа от сервера
swagger:response dicRep
*/
type dicRepBody struct {
	// in:body
	Body struct {
		Hs   string         `json:"hs"`
		Data entities.DicSt `json:"data"`
	}
}
