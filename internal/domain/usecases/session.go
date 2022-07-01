package usecases

import (
	"context"

	"github.com/rendau/dop/dopErrs"
	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
)

func (u *St) SessionGetFromToken(token string) *entities.Session {
	return u.cr.Session.GetFromToken(token)
}

func (u *St) SessionRequireAuth(ses *entities.Session) error {
	if ses.Id == 0 {
		return dopErrs.NotAuthorized
	}

	return nil
}

func (u *St) SessionRequireOneOfRoles(ses *entities.Session, strict bool, roles ...string) error {
	err := u.SessionRequireAuth(ses)
	if err != nil {
		return err
	}

	for _, sRole := range ses.Roles {
		if !strict && sRole == cns.RoleAdmin {
			return nil
		}

		for _, pRole := range roles {
			if pRole == sRole {
				return nil
			}
		}
	}

	return dopErrs.PermissionDenied
}

func (u *St) SessionSetToContext(ctx context.Context, ses *entities.Session) context.Context {
	return u.cr.Session.SetToContext(ctx, ses)
}

func (u *St) SessionSetToContextByToken(ctx context.Context, token string) context.Context {
	return u.cr.Session.SetToContext(ctx, u.SessionGetFromToken(token))
}

func (u *St) SessionGetFromContext(ctx context.Context) *entities.Session {
	return u.cr.Session.GetFromContext(ctx)
}
