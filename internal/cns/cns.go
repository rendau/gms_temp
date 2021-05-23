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

// Static file directories
const (
	SFDUsrAva = "usr_avatar"
)

const (
	UsrTypeUndefined = 0
	UsrTypeAdmin     = 1
)

func UsrTypeIsValid(v int) bool {
	return v == UsrTypeUndefined ||
		v == UsrTypeAdmin
}

// Push types
const (
	NfTypeRefreshProfile = "refresh-profile"
	NfTypeRefreshNumbers = "refresh-numbers"
)
