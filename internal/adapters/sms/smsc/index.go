package smsc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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
		url:        strings.TrimRight(url, "/") + "/",
		httpClient: &http.Client{Timeout: conTimeout},
	}
}

func (s *St) Send(phones string, msg string) bool {
	reqBytes, err := json.Marshal(sendReqSt{
		Phones:  phones,
		Message: msg,
		Sync:    true,
	})
	if err != nil {
		s.lg.Errorw("Fail to marshal json", err)
		return false
	}

	req, err := http.NewRequest("POST", s.url+"send", bytes.NewBuffer(reqBytes))
	if err != nil {
		s.lg.Errorw("Fail to create http-request", err)
		return false
	}

	rep, err := s.httpClient.Do(req)
	if err != nil {
		s.lg.Errorw("Fail to send http-request", err)
		return false
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		s.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return false
	}

	repBytes, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		s.lg.Errorw("Fail to read body", err)
		return false
	}

	if len(repBytes) > 0 {
		var repObj ErrorRepSt

		if err = json.Unmarshal(repBytes, &repObj); err != nil {
			s.lg.Errorw("Fail to parse http-body", err)
			return false
		}

		if repObj.ErrorCode != "" {
			s.lg.Infow("Error from sms", "error_code", repObj.ErrorCode)
			return false
		}
	}

	return true
}

func (s *St) GetBalance() (bool, float64) {
	req, err := http.NewRequest("GET", s.url+"balance", nil)
	if err != nil {
		s.lg.Errorw("Fail to create http-request", err)
		return false, 0
	}

	rep, err := s.httpClient.Do(req)
	if err != nil {
		s.lg.Errorw("Fail to send http-request", err)
		return false, 0
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		s.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return false, 0
	}

	var repObj checkBalanceRepSt

	if err = json.NewDecoder(rep.Body).Decode(&repObj); err != nil {
		s.lg.Errorw("Fail to parse http-body", err)
		return false, 0
	}

	if repObj.ErrorCode != "" {
		s.lg.Infow("Error from sms", "error_code", repObj.ErrorCode)
		return false, 0
	}

	return true, repObj.Balance
}
