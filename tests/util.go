package tests

import (
	"context"
	"testing"

	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
	"github.com/stretchr/testify/require"
)

func resetDb() {
	// var err error
	//
	// bgCtx := context.Background()

	truncateTables([]string{})
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
	// var err error

	app.cache.Clean()

	truncateTables([]string{})

	// _, err = app.db.DbExec(context.Background(), `
	// 	delete from usr where id not in ($1, $2, $3, $4, $5)
	// `, admId, usr1Id, usr2Id, usr3Id, usr4Id)
	// if err != nil {
	// 	app.log.Fatal(err)
	// }
}

func domainErrIsEqual(t *testing.T, v error, expectedErr error, msgArgs ...interface{}) {
	if expectedErr == nil {
		require.Nil(t, v, msgArgs...)
	} else {
		require.NotNil(t, v, msgArgs...)

		switch cErr := v.(type) {
		case errs.Err:
			require.Equal(t, expectedErr.Error(), cErr.Error(), msgArgs...)
		default:
			t.Fatal("bad error type: " + v.Error())
		}
	}
}

func ctxWithSes(ctx context.Context, usrId int64) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return app.ucs.ContextWithSession(ctx, &entities.Session{
		ID: usrId,
	})
}
