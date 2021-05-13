package usecases

import (
	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/core"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
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

func (u *St) SesRequireAuth(ses *entities.Session) error {
	if ses.ID == 0 {
		return errs.NotAuthorized
	}

	return nil
}

func (u *St) SesRequireOneOfTypeIds(ses *entities.Session, typeIds ...int) error {
	err := u.SesRequireAuth(ses)
	if err != nil {
		return err
	}

	if ses.TypeId == cns.UsrTypeAdmin {
		return nil
	}

	for _, typeId := range typeIds {
		if typeId == ses.TypeId {
			return nil
		}
	}

	return errs.PermissionDenied
}
