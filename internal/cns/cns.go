package cns

import "time"

const (
	AppName = "GmsTemp"
	AppUrl  = "https://gmstemp.com"

	MaxPageSize = 1000
)

var (
	AppTimeLocation = time.FixedZone("AST", 21600) // +0600
)

// Roles
const (
	RoleGuest = "guest"
	RoleAdmin = "admin"
)

func RoleIsValid(v string) bool {
	return v == RoleGuest ||
		v == RoleAdmin
}
