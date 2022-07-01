package tests

import (
	"os"
	"testing"

	"github.com/rendau/dop/adapters/cache/mem"
	dopDbPg "github.com/rendau/dop/adapters/db/pg"
	"github.com/rendau/dop/adapters/logger/zap"
	repoPg "github.com/rendau/gms_temp/internal/adapters/repo/pg"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/domain/usecases"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile("conf.yml")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	app.lg = zap.New("info", true)

	app.cache = mem.New()

	app.db, err = dopDbPg.New(true, app.lg, dopDbPg.OptionsSt{
		Dsn: viper.GetString("pg_dsn"),
	})
	if err != nil {
		app.lg.Fatal(err)
	}

	app.repo = repoPg.New(app.db, app.lg)

	app.ucs = usecases.New(app.lg, app.db)

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
		app.repo,
		true,
	)

	app.ucs.SetCore(app.core)

	resetDb()

	// Start tests
	code := m.Run()

	// time.Sleep(20 * time.Second)

	os.Exit(code)
}
