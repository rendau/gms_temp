package usecases

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

const sessionContextKey = "user_session"

func (u *St) ContextWithSession(ctx context.Context, ses *entities.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, ses)
}

func (u *St) ContextGetSession(ctx context.Context) *entities.Session {
	contextV := ctx.Value(sessionContextKey)
	if contextV == nil {
		return &entities.Session{}
	}

	switch ses := contextV.(type) {
	case *entities.Session:
		return ses
	default:
		u.lg.Fatal("wrong type of session in context")
		return &entities.Session{}
	}
}

func (u *St) ContextRequireAuth(ctx context.Context) error {
	return u.SesRequireAuth(u.ContextGetSession(ctx))
}
