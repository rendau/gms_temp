package cfugo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rendau/gms_temp/internal/interfaces"
)

const conTimeout = 10 * time.Second

type St struct {
	lg               interfaces.Logger
	url              string
	apiKey           string
	channelNamespace string

	httpClient *http.Client
}

func New(lg interfaces.Logger, url, apiKey, channelNamespace string) *St {
	return &St{
		lg:               lg,
		url:              strings.TrimRight(url, "/") + "/",
		apiKey:           apiKey,
		channelNamespace: channelNamespace,
		httpClient:       &http.Client{Timeout: conTimeout},
	}
}

func (p *St) Send(channel string, data interface{}) {
	reqObj := sendReqSt{Method: "publish"}
	reqObj.Params.Channel = p.channelNamespace + ":" + channel
	reqObj.Params.Data = data

	dataRaw, err := json.Marshal(reqObj)
	if err != nil {
		p.lg.Errorw("Fail to marshal data", err)
		return
	}

	req, err := http.NewRequest("POST", p.url+"api", bytes.NewBuffer(dataRaw))
	if err != nil {
		p.lg.Errorw("Fail to create http-request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "apikey "+p.apiKey)

	rep, err := p.httpClient.Do(req)
	if err != nil {
		p.lg.Errorw("Fail to send http-request", err)
		return
	}
	defer rep.Body.Close()

	repDataRaw, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		p.lg.Errorw("Fail to read body", err)
		return
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		p.lg.Errorw(
			"Fail to send http-request, bad status code", nil,
			"status_code", rep.StatusCode,
			"body", string(repDataRaw),
		)
		return
	}

	if len(repDataRaw) > 0 {
		repData := sendRepSt{}

		err = json.Unmarshal(repDataRaw, &repData)
		if err != nil {
			p.lg.Errorw(
				"Fail to unmarshal body", err,
				"body", string(repDataRaw),
			)
			return
		}

		if repData.Error != nil {
			p.lg.Errorw(
				"Error in response", nil,
				"body", string(repDataRaw),
			)
			return
		}
	}
}
