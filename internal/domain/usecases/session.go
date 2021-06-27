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

func (u *St) SessionRequireOneOfTypeIds(ses *entities.Session, strict bool, typeIds ...int) error {
	err := u.SessionRequireAuth(ses)
	if err != nil {
		return err
	}

	if !strict {
		if ses.TypeId == cns.UsrTypeAdmin {
			return nil
		}
	}

	for _, typeId := range typeIds {
		if typeId == ses.TypeId {
			return nil
		}
	}

	return errs.PermissionDenied
}

func (u *St) SessionSetToContext(ctx context.Context, ses *entities.Session) context.Context {
	return u.cr.Session.SetToContext(ctx, ses)
}

func (u *St) SessionSetToContextByToken(ctx context.Context, token string) context.Context {
	return u.cr.Session.SetToContext(ctx, u.SessionGet(ctx, token))
}

func (u *St) SessionGetFromContext(ctx context.Context) *entities.Session {
	return u.cr.Session.GetFromContext(ctx)
}
