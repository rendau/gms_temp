package cmd

import (
	"context"
	memCache "github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/cache/redis"
	"github.com/rendau/gms_temp/internal/adapters/httpapi"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/interfaces"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func Execute() {
	var err error

	loadConf()

	app := struct {
		log     *zap.St
		cache   interfaces.Cache
		core    *core.St
		restApi *httpapi.St
	}{}

	app.log, err = zap.New(viper.GetString("log_level"), viper.GetBool("debug"), false)
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetString("redis.url") == "" {
		app.cache = memCache.New()
	} else {
		app.cache = redis.NewRedisSt(
			app.log,
			viper.GetString("redis.url"),
			viper.GetString("redis.psw"),
			viper.GetInt("redis.db"),
		)
	}

	app.core = core.New(
		app.log,
		app.cache,
	)

	app.restApi = httpapi.New(
		app.log,
		viper.GetString("http_listen"),
		app.core,
	)

	app.log.Infow(
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

	app.log.Infow("Shutting down...")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer ctxCancel()

	err = app.restApi.Shutdown(ctx)
	if err != nil {
		app.log.Errorw("Fail to shutdown http-api", err)
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
