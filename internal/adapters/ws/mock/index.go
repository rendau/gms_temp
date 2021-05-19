package mock

import (
	"sync"
)

type St struct {
	q  []Req
	mu sync.Mutex
}

type Req struct {
	UsrIds []int64
	Data   map[string]string
}

func New() *St {
	return &St{
		q: make([]Req, 0),
	}
}

func (m *St) Send2User(usrId int64, data map[string]string) bool {
	return m.Send2Users([]int64{usrId}, data)
}

func (m *St) Send2Users(usrIds []int64, data map[string]string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.q) > 1000 {
		m.q = make([]Req, 0)
	}

	req := Req{
		UsrIds: usrIds,
		Data:   data,
	}

	// fmt.Printf("ws: %+v\n", req)

	m.q = append(m.q, req)

	return true
}

func (m *St) PullAll() []Req {
	m.mu.Lock()
	defer m.mu.Unlock()

	q := m.q

	m.q = make([]Req, 0)

	return q
}

func (m *St) Clean() {
	_ = m.PullAll()
}
