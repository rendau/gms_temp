package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

const sessionContextKey = "user_session"
const sessionDur = int64(600 * time.Second)

type Session struct {
	r *St
}

func NewSession(r *St) *Session {
	return &Session{r: r}
}

func (c *Session) GetFromToken(token string) *entities.Session {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) == 3 {
		if claimsRaw, err := base64.RawURLEncoding.DecodeString(tokenParts[1]); err == nil {
			claims := entities.JwtClaimsSt{}

			err = json.Unmarshal(claimsRaw, &claims)
			if err != nil {
				claims = entities.JwtClaimsSt{}
			}

			claims.Id, _ = strconv.ParseInt(claims.Sub, 10, 64)

			if claims.Roles == nil {
				claims.Roles = make([]string, 0)
			}

			return &claims.Session
		}
	}

	return &entities.Session{}
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
		c.r.lg.Fatal("wrong type of session in context")
		return &entities.Session{}
	}
}

func (c *Session) CreateToken(ses *entities.Session) (string, error) {
	token, _ := c.r.jwts.JwtCreate(
		strconv.FormatInt(ses.Id, 10),
		sessionDur,
		map[string]interface{}{
			"role": ses.Roles,
		},
	)

	return token, nil
}
