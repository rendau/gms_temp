package usecases

import (
	"context"
)

func (u *St) DicGetJson(ctx context.Context,
	reqHs string) (string, []byte, error) {
	return u.cr.Dic.GetJson(ctx, reqHs)
}
