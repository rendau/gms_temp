package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/util"
)

type Notification struct {
	r *St
}

func NewNotification(r *St) *Notification {
	return &Notification{r: r}
}

func (c *Notification) SendRefreshProfile(usrIds []int64) {
	if len(usrIds) == 0 {
		return
	}

	c.r.ws.Send2Users(usrIds, map[string]string{
		"type": cns.NfTypeRefreshProfile,
	})
}

func (c *Notification) SendRefreshNumbers(usrIds []int64) {
	if len(usrIds) == 0 {
		return
	}

	c.r.ws.Send2Users(
		usrIds,
		map[string]string{"type": cns.NfTypeRefreshNumbers},
	)
}

func (c *Notification) SendSmsBalanceAlarm(balance int64) {
	c.r.wg.Add(1)

	if c.r.testing {
		c.sendSmsBalanceAlarmRoutine(balance)
	} else {
		go c.sendSmsBalanceAlarmRoutine(balance)
	}
}

func (c *Notification) sendSmsBalanceAlarmRoutine(balance int64) {
	defer c.r.wg.Done()

	ctx := context.Background()

	admUsrs, _, err := c.r.Usr.List(ctx, &entities.UsrListParsSt{
		// PaginationParams: entities.PaginationParams{Limit: 20},
		TypeId: util.NewInt(cns.UsrTypeAdmin),
	})
	if err != nil {
		return
	}

	usrPhones := make([]string, 0, len(admUsrs))

	for _, usr := range admUsrs {
		usrPhones = append(usrPhones, usr.Phone)
	}

	message := fmt.Sprintf("%s: SMS баланс составляет %dт", cns.AppName, balance)

	c.r.sms.Send(strings.Join(usrPhones, ","), message)
}
