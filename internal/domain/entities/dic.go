package entities

type DicDataSt struct {
	Config *ConfigSt `json:"config"`
	Roles  []*RoleSt `json:"roles"`
}
