package usecases

func (u *St) SystemCronTick5m() {
	u.cr.System.CronTick5m()
}

func (u *St) SystemCronTick15m() {
	u.cr.System.CronTick15m()
}

func (u *St) SystemCronTick30m() {
	u.cr.System.CronTick30m()
}
