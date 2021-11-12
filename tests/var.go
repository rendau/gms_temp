package tests

import (
	"github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/domain/usecases"
)

var (
	app = struct {
		lg    *zap.St
		cache *mem.St
		db    *pg.St
		core  *core.St
		ucs   *usecases.St
	}{}
)
