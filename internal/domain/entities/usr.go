package entities

import (
	"time"
)

type UsrSt struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	TypeId    int       `json:"type_id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
}

type UsrGetPars struct {
	Id    *int64
	Phone *string
	Token *string
}

type UsrProfileSt struct {
	UsrSt
}

type UsrListSt struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	TypeId    int       `json:"type_id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
}

type UsrListParsSt struct {
	PaginationParams

	Ids    *[]int64
	TypeId *int
	Search *string

	SortBy *string //

	OnlyCount bool
}

type UsrCUSt struct {
	TypeId *int    `json:"type_id"`
	Phone  *string `json:"-"`
	Name   *string `json:"name"`
}

type UsrAuthReqSt struct {
	Phone   string `json:"phone"`
	SmsCode int    `json:"sms_code"`
}

type UsrRegisterSt struct {
	Phone   string `json:"phone"`
	SmsCode int    `json:"sms_code"`

	TypeId *int    `json:"type_id"`
	Name   *string `json:"name"`
}
