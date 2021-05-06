package websocket

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rendau/gms_temp/internal/interfaces"
)

const conTimeout = 20 * time.Second

type St struct {
	lg  interfaces.Logger
	url string

	httpClient *http.Client
}

func New(lg interfaces.Logger, url string) *St {
	return &St{
		lg:         lg,
		url:        url,
		httpClient: &http.Client{Timeout: conTimeout},
	}
}

func (p *St) Send2User(usrId int64, data map[string]string) bool {
	return p.Send2Users([]int64{usrId}, data)
}

func (p *St) Send2Users(usrIds []int64, data map[string]string) bool {
	if len(usrIds) == 0 {
		return true
	}

	message, err := json.Marshal(data)
	if err != nil {
		p.lg.Errorw("Fail to marshal data", err)
		return false
	}

	reqBytes, err := json.Marshal(sendReqSt{
		UsrIds:  usrIds,
		Message: message,
	})
	if err != nil {
		p.lg.Errorw("Fail to marshal json", err)
		return false
	}

	req, err := http.NewRequest("POST", p.url+"send", bytes.NewBuffer(reqBytes))
	if err != nil {
		p.lg.Errorw("Fail to create http-request", err)
		return false
	}

	rep, err := p.httpClient.Do(req)
	if err != nil {
		p.lg.Errorw("Fail to send http-request", err)
		return false
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		p.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return false
	}

	return true
}
