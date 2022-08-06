package core

import (
	"context"
	"strconv"

	"github.com/rendau/dop/adapters/jwt"
	"github.com/rendau/gms_temp/internal/domain/entities"
)

const sessionContextKey = "user_session"

type Session struct {
	r *St
}

func NewSession(r *St) *Session {
	return &Session{r: r}
}

func (c *Session) GetFromToken(token string) *entities.Session {
	var session entities.Session

	if jwt.ParsePayload(token, &session) != nil {
		session = entities.Session{}
	}

	session.Id, _ = strconv.ParseInt(session.Sub, 10, 64)

	if session.Roles == nil {
		session.Roles = make([]string, 0)
	}

	return &session
}

func (c *Session) SetToContext(ctx context.Context, ses *entities.Session) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, sessionContextKey, ses)
}

func (c *Session) GetFromContext(ctx context.Context) *entities.Session {
	contextV := ctx.Value(sessionContextKey)
	if contextV == nil {
		return &entities.Session{}
	}

	switch ses := contextV.(type) {
	case *entities.Session:
		return ses
	default:
		c.r.lg.Errorw("wrong type of session in context", nil)
		return &entities.Session{}
	}
}
