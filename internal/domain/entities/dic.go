package entities

type DicSt struct {
	UsrTypes []*UsrTypeSt `json:"usr_types"`
	Config   *ConfigSt    `json:"config"`
}
