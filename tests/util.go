package tests

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

func resetDb() {
	truncateTables([]string{
		"cfg",
	})
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

	truncateTables([]string{
		"cfg",
	})

	// _, err = app.db.DbExec(context.Background(), `
	// 	delete from tablename where x = $1
	// `, 111)
	// if err != nil {
	// 	app.lg.Fatal(err)
	// }
}

func ctxWithSes(ses *entities.Session) context.Context {
	return app.ucs.SessionSetToContext(nil, ses)
}
