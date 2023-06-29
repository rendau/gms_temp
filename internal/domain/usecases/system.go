package usecases

import (
	"time"
)

func (u *St) SystemCronTick5m(t time.Time) {
	u.cr.System.CronTick5m(t)
}

func (u *St) SystemCronTick15m(t time.Time) {
	u.cr.System.CronTick15m(t)
}

func (u *St) SystemCronTick30m(t time.Time) {
	u.cr.System.CronTick30m(t)
}
