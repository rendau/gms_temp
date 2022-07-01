package repo

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

type Repo interface {
	// config
	ConfigGet(ctx context.Context) (*entities.ConfigSt, error)
	ConfigSet(ctx context.Context, config *entities.ConfigSt) error
}
