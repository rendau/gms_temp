package cns

import "time"

const (
	AppName = "GmsTemp"
	AppUrl  = "https://gmstemp.kz"

	MaxPageSize = 1000
)

var (
	AppTimeLocation = time.FixedZone("AST", 21600) // +0600
)

const (
	UsrTypeUndefined = 0
	UsrTypeAdmin     = 1
)

func UsrTypeIsValid(v int) bool {
	return v == UsrTypeUndefined ||
		v == UsrTypeAdmin
}
