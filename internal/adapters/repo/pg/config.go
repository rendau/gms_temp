package pg

import (
	"context"
	"errors"

	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/dopErrs"
	"github.com/rendau/gms_temp/internal/domain/entities"
)

type ConfigRowSt struct {
	V entities.ConfigSt `json:"v" db:"v"`
}

func (d *St) ConfigGet(ctx context.Context) (*entities.ConfigSt, error) {
	item := ConfigRowSt{}

	err := d.HfGet(ctx, db.RDBGetOptions{
		Dst:    &item,
		Tables: []string{"cfg"},
	})
	if err != nil {
		if errors.Is(err, dopErrs.NoRows) {
			return &item.V, nil
		}
		return nil, err
	}

	return &item.V, nil
}

func (d *St) ConfigSet(ctx context.Context, config *entities.ConfigSt) error {
	err := d.DbExec(ctx, `
		with u as (
			update cfg
			set v = $1
			returning 1
		)
		insert into cfg (v)
		select $1
		where not exists(select * from u)
	`, config)
	if err != nil {
		return err
	}

	return nil
}
