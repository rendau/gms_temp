package entities

type DicDataSt struct {
	Config   *ConfigSt    `json:"config"`
	UsrTypes []*UsrTypeSt `json:"usr_types"`
}
