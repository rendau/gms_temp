package entities

import (
	"time"
)

type UsrSt struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	TypeId    int       `json:"type_id"`
	Phone     string    `json:"phone"`
	Ava       string    `json:"ava"`
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

type UsrNumbersSt struct {
	// NewMsgCount int64 `json:"new_msg_count"`
}

type UsrListSt struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	TypeId    int       `json:"type_id"`
	Phone     string    `json:"phone"`
	Ava       string    `json:"ava"`
	Name      string    `json:"name"`
}

type UsrListParsSt struct {
	PaginationParams

	Ids     *[]int64
	TypeId  *int
	TypeIds *[]int
	Search  *string

	SortBy *string //

	OnlyCount bool
}

type UsrCUSt struct {
	TypeId *int    `json:"type_id"`
	Phone  *string `json:"-"`
	Name   *string `json:"name"`
	Ava    *string `json:"ava"`
}

type PhoneAndSmsCodeSt struct {
	Phone   string `json:"phone"`
	SmsCode int    `json:"sms_code"`
}

type UsrRegReqSt struct {
	PhoneAndSmsCodeSt

	TypeId *int    `json:"type_id"`
	Name   *string `json:"name"`
	Ava    *string `json:"ava"`
}
