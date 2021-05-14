package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

func (u *St) ProfileSendPhoneValidatingCode(ctx context.Context,
	phone string, errNE bool) error {
	var err error

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	err = u.cr.Usr.SendPhoneValidatingCode(ctx, phone, errNE)
	if err != nil {
		return err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func (u *St) ProfileAuth(ctx context.Context,
	obj *entities.UsrAuthReqSt) (int64, string, error) {
	var err error

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return 0, "", err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	id, token, err := u.cr.Usr.Auth(ctx, obj)
	if err != nil {
		return 0, "", err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return 0, "", err
	}

	return id, token, nil
}

func (u *St) ProfileLogout(ctx context.Context) error {
	var err error

	ses := u.ContextGetSession(ctx)

	if ses.ID == 0 {
		return nil
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	err = u.cr.Usr.Logout(ctx, ses.ID)
	if err != nil {
		return err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func (u *St) ProfileGet(ctx context.Context) (*entities.UsrProfileSt, error) {
	var err error

	ses := u.ContextGetSession(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return nil, err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	result, err := u.cr.Usr.GetProfile(ctx, ses.ID)
	if err != nil {
		return nil, err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

func (u *St) ProfileUpdate(ctx context.Context,
	obj *entities.UsrCUSt) error {
	var err error

	ses := u.ContextGetSession(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	// restrict
	obj.TypeId = nil
	obj.Phone = nil

	err = u.cr.Usr.Update(ctx, ses.ID, obj)
	if err != nil {
		return err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func (u *St) ProfileGetId(ctx context.Context) (int64, error) {
	return u.ContextGetSession(ctx).ID, nil
}

func (u *St) ProfileDelete(ctx context.Context) error {
	var err error

	ses := u.ContextGetSession(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	err = u.cr.Usr.Delete(ctx, ses.ID)
	if err != nil {
		return err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return err
	}

	return nil
}
