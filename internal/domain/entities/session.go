package entities

type Session struct {
	Id    int64    `json:"id"`
	Roles []string `json:"roles"`
}

type JwtClaimsSt struct {
	Session

	Sub string `json:"sub"`
}
