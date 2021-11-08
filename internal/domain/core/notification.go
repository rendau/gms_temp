package core

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/util"
)

const NotificationWsChannelPrefixUsr = "usr#"

type Notification struct {
	r *St
}

func NewNotification(r *St) *Notification {
	return &Notification{r: r}
}

func (c *Notification) SendRefreshProfile(usrId int64) {
	c.r.ws.Send(NotificationWsChannelPrefixUsr+strconv.FormatInt(usrId, 10), map[string]string{
		"type": cns.NfTypeRefreshProfile,
	})
}

func (c *Notification) SendRefreshProfileMany(usrIds []int64) {
	if len(usrIds) == 0 {
		return
	}

	for _, id := range usrIds {
		c.SendRefreshProfile(id)
	}
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
