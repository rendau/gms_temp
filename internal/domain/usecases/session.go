package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
)

func (u *St) SessionGet(ctx context.Context, token string) *entities.Session {
	return u.cr.Session.Get(ctx, token)
}

func (u *St) SessionRequireAuth(ses *entities.Session) error {
	if ses.ID == 0 {
		return errs.NotAuthorized
	}

	return nil
}

func (u *St) SessionRequireOneOfTypeIds(ses *entities.Session, typeIds ...int) error {
	err := u.SessionRequireAuth(ses)
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
