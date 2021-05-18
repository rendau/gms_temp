package interfaces

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

type Db interface {
	ContextWithTransaction(ctx context.Context) (context.Context, error)
	CommitContextTransaction(ctx context.Context) error
	RollbackContextTransaction(ctx context.Context)
	RenewContextTransaction(ctx context.Context) error

	// config
	ConfigGet(ctx context.Context) (*entities.ConfigSt, error)
	ConfigSet(ctx context.Context, config *entities.ConfigSt) error

	// usr
	UsrList(ctx context.Context, pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error)
	UsrGet(ctx context.Context, pars *entities.UsrGetPars) (*entities.UsrSt, error)
	UsrIdExists(ctx context.Context, id int64) (bool, error)
	UsrIdsExists(ctx context.Context, ids []int64) (bool, error)
	UsrPhoneExists(ctx context.Context, phone string, excludeId int64) (bool, error)
	UsrGetToken(ctx context.Context, id int64) (string, error)
	UsrSetToken(ctx context.Context, id int64, token string) error
	UsrGetIdForToken(ctx context.Context, token string) (int64, error)
	UsrGetTypeId(ctx context.Context, id int64) (int, error)
	UsrGetPhone(ctx context.Context, id int64) (string, error)
	UsrGetIdForPhone(ctx context.Context, phone string) (int64, error)
	UsrCreate(ctx context.Context, obj *entities.UsrCUSt) (int64, error)
	UsrUpdate(ctx context.Context, id int64, obj *entities.UsrCUSt) error
	UsrDelete(ctx context.Context, id int64) error
}
