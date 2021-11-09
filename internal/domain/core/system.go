package core

type System struct {
	r *St
}

func NewSystem(r *St) *System {
	return &System{r: r}
}

func (c *System) CronTick5m() {
	// do something in goroutine
}

func (c *System) CronTick15m() {
	// do something in goroutine
}

func (c *System) CronTick30m() {
	// do something in goroutine
}
