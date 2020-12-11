package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rendau/gms_temp/internal/domain/usecases"

	memCache "github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/cache/redis"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	"github.com/rendau/gms_temp/internal/adapters/httpapi/rest"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/interfaces"
	"github.com/spf13/viper"
)

func Execute() {
	var err error

	loadConf()

	app := struct {
		lg      *zap.St
		cache   interfaces.Cache
		db      interfaces.Db
		core    *core.St
		ucs     *usecases.St
		restApi *rest.St
	}{}

	app.lg, err = zap.New(viper.GetString("log_level"), viper.GetBool("debug"), false)
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetString("redis.url") == "" {
		app.cache = memCache.New()
	} else {
		app.cache = redis.NewRedisSt(
			app.lg,
			viper.GetString("redis.url"),
			viper.GetString("redis.psw"),
			viper.GetInt("redis.db"),
		)
	}

	app.db, err = pg.NewSt(app.lg, viper.GetString("pg.dsn"))
	if err != nil {
		app.lg.Fatal(err)
	}

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
	)

	app.ucs = usecases.New(
		app.lg,
		app.db,
		app.cache,
		app.core,
	)

	app.restApi = rest.New(
		app.lg,
		viper.GetString("http_listen"),
		app.ucs,
	)

	app.lg.Infow(
		"Starting",
		"http_listen", viper.GetString("http_listen"),
	)

	app.restApi.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var exitCode int

	select {
	case <-stop:
	case <-app.restApi.Wait():
		exitCode = 1
	}

	app.lg.Infow("Shutting down...")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer ctxCancel()

	err = app.restApi.Shutdown(ctx)
	if err != nil {
		app.lg.Errorw("Fail to shutdown http-api", err)
		exitCode = 1
	}

	os.Exit(exitCode)
}

func loadConf() {
	viper.SetDefault("debug", "false")
	viper.SetDefault("http_listen", ":80")
	viper.SetDefault("log_level", "debug")

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	// viper.Set("some.url", uriRPadSlash(viper.GetString("some.url")))
}

func uriRPadSlash(uri string) string {
	if uri != "" && !strings.HasSuffix(uri, "/") {
		return uri + "/"
	}
	return uri
}
