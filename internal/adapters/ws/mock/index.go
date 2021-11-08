package mock

import (
	"sync"

	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg      interfaces.Logger
	testing bool

	q  []Req
	mu sync.Mutex
}

type Req struct {
	channel string
	Data    interface{}
}

func New(lg interfaces.Logger, testing bool) *St {
	return &St{
		lg:      lg,
		testing: testing,
		q:       make([]Req, 0),
	}
}

func (m *St) Send(channel string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.q) > 1000 {
		m.q = make([]Req, 0)
	}

	req := Req{
		channel: channel,
		Data:    data,
	}

	if !m.testing {
		m.lg.Infow("Ws sent", "channel", channel, "data", data)
	}

	m.q = append(m.q, req)
}

func (m *St) PullAll() []Req {
	m.mu.Lock()
	defer m.mu.Unlock()

	q := m.q

	m.q = make([]Req, 0)

	return q
}

func (m *St) Get(channel string) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, req := range m.q {
		if req.channel == channel {
			return req.Data
		}
	}

	return nil
}

func (m *St) Clean() {
	_ = m.PullAll()
}
