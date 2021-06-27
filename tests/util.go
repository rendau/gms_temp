package tests

import (
	"context"
	"testing"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

func resetDb() {
	var err error

	truncateTables([]string{
		"cfg", "usr",
	})

	bgCtx := context.Background()

	usrs := []struct {
		IdPtr  *int64
		Name   string
		Phone  string
		TypeId int
	}{
		{&admId, admName, admPhone, cns.UsrTypeAdmin},
		{&usr1Id, usr1Name, usr1Phone, cns.UsrTypeUndefined},
	}
	for _, usr := range usrs {
		*usr.IdPtr, err = app.core.Usr.Create(bgCtx, &entities.UsrCUSt{
			TypeId: &usr.TypeId,
			Name:   &usr.Name,
			Phone:  &usr.Phone,
		})
		if err != nil {
			app.lg.Fatal(err)
		}
	}
}

func truncateTables(tables []string) {
	q := ``
	for _, t := range tables {
		q += ` truncate ` + t + ` restart identity cascade; `
	}
	if q != `` {
		_, err := app.db.DbExec(context.Background(), `begin; `+q+` commit;`)
		if err != nil {
			app.lg.Fatal(err)
		}
	}
}

func prepareDbForNewTest() {
	var err error

	app.cache.Clean()
	app.sms.Clean()
	app.ws.Clean()

	truncateTables([]string{
		"cfg",
	})

	_, err = app.db.DbExec(context.Background(), `
		delete from usr where id not in (select * from unnest($1 :: bigint[]))
	`, []int64{admId, usr1Id})
	if err != nil {
		app.lg.Fatal(err)
	}
}

func ctxWithSes(t *testing.T, ctx context.Context, usrId int64) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	token, err := app.core.Usr.GetOrCreateToken(ctx, usrId)
	require.Nil(t, err)

	return app.ucs.SessionSetToContextByToken(ctx, token)
}
