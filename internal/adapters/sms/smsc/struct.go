package smsc

type sendReqSt struct {
	Phones  string `json:"phones"`
	Message string `json:"message"`
	Sync    bool   `json:"sync"`
}

type ErrorRepSt struct {
	ErrorCode string `json:"error_code"`
}

type checkBalanceRepSt struct {
	ErrorRepSt
	Balance float64 `json:"balance"`
}
