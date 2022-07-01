package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

func (u *St) DicGet(ctx context.Context) (*entities.DicSt, error) {
	return u.cr.Dic.Get(ctx)
}
