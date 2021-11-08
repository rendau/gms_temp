package cfugo

type sendReqSt struct {
	Method string `json:"method"`
	Params struct {
		Channel string      `json:"channel"`
		Data    interface{} `json:"data"`
	} `json:"params"`
}

type sendRepSt struct {
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
