package websocket

import (
	"encoding/json"
)

type sendReqSt struct {
	UsrIds  []int64         `json:"usr_ids"`
	Message json.RawMessage `json:"message"`
}
