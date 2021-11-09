package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	memCache "github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/cache/redis"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	"github.com/rendau/gms_temp/internal/adapters/httpapi/rest"
	"github.com/rendau/gms_temp/internal/adapters/jwts/jwts"
	jwtsMock "github.com/rendau/gms_temp/internal/adapters/jwts/mock"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/domain/usecases"
	"github.com/rendau/gms_temp/internal/interfaces"
	"github.com/spf13/viper"
)

func Execute() {
	var err error

	loadConf()

	debug := viper.GetBool("DEBUG")

	app := struct {
		lg      *zap.St
		cache   interfaces.Cache
		db      interfaces.Db
		jwts    interfaces.Jwts
		core    *core.St
		ucs     *usecases.St
		restApi *rest.St
	}{}

	app.lg, err = zap.New(viper.GetString("LOG_LEVEL"), debug, false)
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetString("REDIS_URL") == "" {
		app.cache = memCache.New()
	} else {
		app.cache = redis.New(
			app.lg,
			viper.GetString("REDIS_URL"),
			viper.GetString("REDIS_PSW"),
			viper.GetInt("REDIS_DB"),
		)
	}

	app.db, err = pg.New(app.lg, viper.GetString("PG_DSN"), debug)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.ucs = usecases.New(
		app.lg,
		app.db,
	)

	if viper.GetString("MS_JWTS_URL") == "" {
		app.jwts = jwtsMock.New(app.lg, false)
	} else {
		app.jwts = jwts.New(
			app.lg,
			viper.GetString("MS_JWTS_URL"),
		)
	}

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
		app.jwts,
		false,
	)

	app.ucs.SetCore(app.core)

	restApiEChan := make(chan error, 1)

	app.restApi = rest.New(
		app.lg,
		viper.GetString("HTTP_LISTEN"),
		app.ucs,
		restApiEChan,
		viper.GetBool("HTTP_CORS"),
	)

	app.lg.Infow("Starting")

	app.restApi.Start()

	stopSignalChan := make(chan os.Signal, 1)
	signal.Notify(stopSignalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var exitCode int

	select {
	case <-stopSignalChan:
	case <-restApiEChan:
		exitCode = 1
	}

	app.lg.Infow("Shutting down...")

	err = app.restApi.Shutdown(20 * time.Second)
	if err != nil {
		app.lg.Errorw("Fail to shutdown http-api", err)
		exitCode = 1
	}

	app.lg.Infow("Wait routines...")

	app.core.WaitJobs()

	app.lg.Infow("Exit")

	os.Exit(exitCode)
}

func loadConf() {
	viper.SetDefault("debug", "false")
	viper.SetDefault("http_listen", ":80")
	viper.SetDefault("log_level", "info")

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()
}
