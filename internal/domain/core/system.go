package core

import (
	"time"
)

type System struct {
	r *St
}

func NewSystem(r *St) *System {
	return &System{r: r}
}

func (c *System) CronTick5m(t time.Time) {
	//t = t.In(cns.AppTimeLocation)
	//hour := t.Hour()
	//minute := t.Minute()

	// do something
}

func (c *System) CronTick15m(t time.Time) {
	//t = t.In(cns.AppTimeLocation)
	//hour := t.Hour()
	//minute := t.Minute()

	// do something
}

func (c *System) CronTick30m(t time.Time) {
	//t = t.In(cns.AppTimeLocation)
	//hour := t.Hour()
	//minute := t.Minute()

	// do something
}
