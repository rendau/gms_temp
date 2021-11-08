package usecases

import (
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg interfaces.Logger

	db interfaces.Db
	cr *core.St
}

func New(
	lg interfaces.Logger,
	db interfaces.Db,
) *St {
	u := &St{
		lg: lg,
		db: db,
	}

	return u
}

func (u *St) SetCore(core *core.St) {
	u.cr = core
}
