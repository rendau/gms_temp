package main

import (
	"github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

const confPath = "test_conf.yml"

var (
	app = struct {
		lg    *zap.St
		cache *mem.St
		core  *core.St
		db    *pg.St
	}{}
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile(confPath)
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

	app.db, err = pg.NewSt(app.lg, viper.GetString("pg.dsn"))
	if err != nil {
		app.lg.Fatal(err)
	}

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
	)

	// Start tests
	code := m.Run()

	os.Exit(code)
}

func TestMenu(t *testing.T) {
	require.True(t, true, true)
}
