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
	cr *core.St,
) *St {
	u := &St{
		lg: lg,
		db: db,
		cr: cr,
	}

	return u
}
