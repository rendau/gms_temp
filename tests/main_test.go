package tests

import (
	"log"
	"os"
	"testing"

	"github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/domain/usecases"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile("conf.yml")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	app.lg, err = zap.New(
		"info",
		true,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer app.lg.Sync()

	app.cache = mem.New()

	app.db, err = pg.New(app.lg, viper.GetString("pg_dsn"), true)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
		true,
	)

	app.ucs = usecases.New(
		app.lg,
		app.db,
		app.core,
	)

	resetDb()

	app.core.Start()

	// Start tests
	code := m.Run()

	os.Exit(code)
}
