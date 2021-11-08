package entities

type Session struct {
	Id     int64 `json:"id"`
	TypeId int   `json:"type_id"`
}

type JwtClaimsSt struct {
	Session

	Sub string `json:"sub"`
}
