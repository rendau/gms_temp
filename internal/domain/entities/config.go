package entities

type ConfigSt struct {
	Contacts ConfigContactsSt `json:"contacts"`
}

type ConfigContactsSt struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}
