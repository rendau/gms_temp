package jwts

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rendau/gms_temp/internal/domain/errs"
	"github.com/rendau/gms_temp/internal/interfaces"
)

const conTimeout = 10 * time.Second

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

func (p *St) JwtCreate(sub string, expSeconds int64, payload map[string]interface{}) (string, error) {
	data := map[string]interface{}{}

	for k, v := range payload {
		data[k] = v
	}

	data["sub"] = sub

	if expSeconds != 0 {
		data["exp_seconds"] = expSeconds
	}

	dataRaw, err := json.Marshal(data)
	if err != nil {
		p.lg.Errorw("Fail to marshal data", err)
		return "", err
	}

	req, err := http.NewRequest("POST", p.url+"jwt", bytes.NewBuffer(dataRaw))
	if err != nil {
		p.lg.Errorw("Fail to create http-request", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rep, err := p.httpClient.Do(req)
	if err != nil {
		p.lg.Errorw("Fail to send http-request", err)
		return "", err
	}
	defer rep.Body.Close()

	repDataRaw, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		p.lg.Errorw("Fail to read body", err)
		return "", err
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		p.lg.Errorw(
			"Fail to send http-request, bad status code", nil,
			"status_code", rep.StatusCode,
			"rep_body", string(repDataRaw),
		)
		return "", errs.ServiceNA
	}

	repData := jwtCreateRepSt{}

	err = json.Unmarshal(repDataRaw, &repData)
	if err != nil {
		p.lg.Errorw(
			"Fail to unmarshal body", err,
			"rep_body", string(repDataRaw),
		)
		return "", err
	}

	return repData.Token, nil
}
