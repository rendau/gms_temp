package mock

import (
	"regexp"
	"strconv"
	"sync"

	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg       interfaces.Logger
	traceMsg bool

	q             []Req
	mu            sync.Mutex
	smsCodeRegexp *regexp.Regexp
}

type Req struct {
	Phones string
	Msg    string
}

func New(lg interfaces.Logger, traceMsg bool) *St {
	return &St{
		lg:       lg,
		traceMsg: traceMsg,

		q:             make([]Req, 0),
		smsCodeRegexp: regexp.MustCompile(`([0-9]{4})`),
	}
}

func (m *St) Send(phones string, msg string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.q) > 100 {
		m.q = make([]Req, 0)
	}

	req := Req{
		Phones: phones,
		Msg:    msg,
	}

	if m.traceMsg {
		m.lg.Infow("Sms", "phones", phones, "msg", msg)
	}

	m.q = append(m.q, req)

	return true
}

func (m *St) GetBalance() (bool, float64) {
	return true, 0
}

func (m *St) PullAll() []Req {
	m.mu.Lock()
	defer m.mu.Unlock()

	q := m.q

	m.q = make([]Req, 0)

	return q
}

func (m *St) PullCode() int {
	smsReqs := m.PullAll()
	if len(smsReqs) < 1 {
		return 0
	}

	matches := m.smsCodeRegexp.FindStringSubmatch(smsReqs[0].Msg)
	if len(matches) == 2 {
		code, _ := strconv.ParseInt(matches[1], 10, 64)
		return int(code)
	}

	return 0
}

func (m *St) Clean() {
	_ = m.PullAll()
}
