package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/util"
)

func (u *St) UsrList(ctx context.Context,
	pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireOneOfTypeIds(ses, false); err != nil {
		return nil, 0, err
	}

	if err = util.RequirePageSize(pars.PaginationParams, 0); err != nil {
		return nil, 0, err
	}

	return u.cr.Usr.List(ctx, pars)
}

func (u *St) UsrGet(ctx context.Context,
	id int64) (*entities.UsrSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Usr.Get(ctx, &entities.UsrGetParsSt{Id: &id}, true)
}

func (u *St) UsrCreate(ctx context.Context,
	obj *entities.UsrCUSt) (int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireOneOfTypeIds(ses, false); err != nil {
		return 0, err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return 0, err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	newId, err := u.cr.Usr.Create(ctx, obj)
	if err != nil {
		return 0, err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return 0, err
	}

	return newId, nil
}

func (u *St) UsrUpdate(ctx context.Context,
	id int64, obj *entities.UsrCUSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireOneOfTypeIds(ses, false); err != nil {
		return err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	err = u.cr.Usr.Update(ctx, id, obj)
	if err != nil {
		return err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func (u *St) UsrDelete(ctx context.Context,
	id int64) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireOneOfTypeIds(ses, false); err != nil {
		return err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	err = u.cr.Usr.Delete(ctx, id)
	if err != nil {
		return err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return err
	}

	return nil
}
