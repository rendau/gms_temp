package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
)

func (u *St) ConfigSet(ctx context.Context,
	config *entities.ConfigSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireOneOfRoles(ses, false, cns.RoleAdmin); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Config.Set(ctx, config)
	})
}

func (u *St) ConfigGet(ctx context.Context) (*entities.ConfigSt, error) {
	return u.cr.Config.Get(ctx)
}
