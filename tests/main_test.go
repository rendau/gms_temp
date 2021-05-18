package tests

import (
	"log"
	"os"
	"testing"

	"github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	smsMock "github.com/rendau/gms_temp/internal/adapters/sms/mock"
	wsMock "github.com/rendau/gms_temp/internal/adapters/ws/mock"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/domain/usecases"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile("test_conf.yml")
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

	app.db, err = pg.New(app.lg, viper.GetString("pg.dsn"), true)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.sms = smsMock.New(false)

	app.ws = wsMock.New()

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
		app.sms,
		app.ws,
		false,
		true,
	)

	app.ucs = usecases.New(
		app.lg,
		app.db,
		app.core,
	)

	resetDb()

	// Start tests
	code := m.Run()

	os.Exit(code)
}
