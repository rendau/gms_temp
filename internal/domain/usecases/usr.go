package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
	"github.com/rendau/gms_temp/internal/domain/util"
)

func (u *St) UsrList(token string,
	pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error) {
	var err error

	ctx := context.Background()

	ses := u.cr.Session.Get(ctx, token)

	if err = u.SesRequireOneOfTypeIds(ses); err != nil {
		return nil, 0, err
	}

	if err = util.RequirePageSize(pars.PaginationParams, 0); err != nil {
		return nil, 0, err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return nil, 0, err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	usrs, totalCnt, err := u.cr.Usr.List(ctx, pars)
	if err != nil {
		return nil, 0, err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return nil, 0, err
	}

	return usrs, totalCnt, nil
}

func (u *St) UsrGet(token string,
	id int64) (*entities.UsrSt, error) {
	var err error

	ctx := context.Background()

	ses := u.cr.Session.Get(ctx, token)

	if err = u.SesRequireAuth(ses); err != nil {
		return nil, err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return nil, err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	usr, err := u.cr.Usr.Get(ctx, &entities.UsrGetPars{Id: &id}, true)
	if err != nil {
		return nil, err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return nil, err
	}

	return usr, nil
}

func (u *St) UsrCreate(token string,
	obj *entities.UsrCUSt) (int64, error) {
	var err error

	ctx := context.Background()

	ses := u.cr.Session.Get(ctx, token)

	if err = u.SesRequireOneOfTypeIds(ses); err != nil {
		return 0, err
	}

	if ses.TypeId != cns.UsrTypeAdmin {
		return 0, errs.PermissionDenied
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

func (u *St) UsrUpdate(token string,
	id int64, obj *entities.UsrCUSt) error {
	var err error

	ctx := context.Background()

	ses := u.cr.Session.Get(ctx, token)

	if err = u.SesRequireOneOfTypeIds(ses); err != nil {
		return err
	}

	if ses.TypeId != cns.UsrTypeAdmin {
		return errs.PermissionDenied
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

func (u *St) UsrDelete(token string,
	id int64) error {
	var err error

	ctx := context.Background()

	ses := u.cr.Session.Get(ctx, token)

	if err = u.SesRequireOneOfTypeIds(ses); err != nil {
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
