package entities

type DicDataSt struct {
	UsrTypes []*UsrTypeSt `json:"usr_types"`
	Config   *ConfigSt    `json:"config"`
}
