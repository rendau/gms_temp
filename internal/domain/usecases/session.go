package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
)

func (u *St) SessionGet(ctx context.Context, token string) *entities.Session {
	return u.cr.Session.Get(ctx, token)
}

func (u *St) SesRequireAuth(ses *entities.Session) error {
	if ses.ID == 0 {
		return errs.NotAuthorized
	}
	return nil
}
