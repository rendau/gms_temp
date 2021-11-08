package mock

import (
	"encoding/base64"
	"encoding/json"

	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg      interfaces.Logger
	testing bool
}

func New(lg interfaces.Logger, testing bool) *St {
	return &St{
		lg:      lg,
		testing: testing,
	}
}

func (p *St) JwtCreate(sub string, expSeconds int64, payload map[string]interface{}) (string, error) {
	payload["sub"] = sub

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		p.lg.Errorw("Fail to marshal data", err)
		return "", err
	}

	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadRaw)

	return "XXX." + payloadB64 + ".YYY", nil
}
