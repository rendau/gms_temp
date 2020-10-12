package main

import (
	"github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

var (
	app = struct {
		log   *zap.St
		cache *mem.St
		core  *core.St
	}{}
)

func TestMain(m *testing.M) {
	var err error

	app.log, err = zap.New(
		"info",
		true,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer app.log.Sync()

	app.cache = mem.New()

	app.core = core.New(
		app.log,
		app.cache,
	)

	// Start tests
	code := m.Run()

	os.Exit(code)
}

func TestMenu(t *testing.T) {
	require.True(t, true, true)
}
