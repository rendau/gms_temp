package tests

import (
	"github.com/rendau/gms_temp/internal/adapters/cache/mem"
	"github.com/rendau/gms_temp/internal/adapters/db/pg"
	jwtsMock "github.com/rendau/gms_temp/internal/adapters/jwts/mock"
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
		jwts  *jwtsMock.St
		sms   *smsMock.St
		ws    *wsMock.St
		core  *core.St
		ucs   *usecases.St
	}{}

	admId    int64
	admName  = "Admin"
	admPhone = "70000000001"

	usr1Id    int64
	usr1Name  = "Usr1"
	usr1Phone = "75550000001"
)
