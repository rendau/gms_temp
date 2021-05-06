package tests

import (
	"github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	"github.com/rendau/gms_temp/internal/adapters/logger/zap"
	smsMock "github.com/rendau/gms_temp/internal/adapters/sms/mock"
	wsMock "github.com/rendau/gms_temp/internal/adapters/ws/mock"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/domain/usecases"
)

var (
	app = struct {
		lg    *zap.St
		cache *mem.St
		db    *pg.St
		sms   *smsMock.St
		ws    *wsMock.St
		core  *core.St
		ucs   *usecases.St
	}{}
)
