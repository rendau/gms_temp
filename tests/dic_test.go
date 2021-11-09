package tests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
	"github.com/stretchr/testify/require"
)

func TestDic(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()
	admCtx := app.ucs.SessionSetToContext(context.Background(), &entities.Session{
		Id:    1,
		Roles: []string{cns.RoleAdmin},
	})

	dicHs, dicJson, err := app.ucs.DicGetJson(bgCtx, "")
	require.Nil(t, err)
	require.NotEmpty(t, dicHs)
	require.NotNil(t, dicJson)
	require.NotEmpty(t, dicJson)

	var dicHs1 string

	dicHs1, dicJson, err = app.ucs.DicGetJson(bgCtx, dicHs)
	require.Nil(t, err)
	require.Equal(t, dicHs, dicHs1)
	require.Nil(t, dicJson)

	dicHs1, dicJson, err = app.ucs.DicGetJson(bgCtx, "")
	require.Nil(t, err)
	require.Equal(t, dicHs, dicHs1)
	require.NotNil(t, dicJson)
	require.NotEmpty(t, dicJson)

	err = app.ucs.ConfigSet(admCtx, &entities.ConfigSt{
		Contacts: entities.ConfigContactsSt{
			Phone: "71230000321",
			Email: "qwe@asd.com",
		},
	})
	require.Nil(t, err)

	dicHs1, dicJson, err = app.ucs.DicGetJson(bgCtx, dicHs)
	require.Nil(t, err)
	require.NotEmpty(t, dicHs)
	require.NotEqual(t, dicHs, dicHs1)
	require.NotNil(t, dicJson)
	require.NotEmpty(t, dicJson)

	dic := &entities.DicDataSt{}

	err = json.Unmarshal(dicJson, dic)
	require.Nil(t, err)
	require.Equal(t, "71230000321", dic.Config.Contacts.Phone)
	require.Equal(t, "qwe@asd.com", dic.Config.Contacts.Email)
}

func TestCfg(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()
	guestCtx := app.ucs.SessionSetToContext(context.Background(), &entities.Session{
		Id:    1,
		Roles: []string{cns.RoleGuest},
	})
	admCtx := app.ucs.SessionSetToContext(context.Background(), &entities.Session{
		Id:    1,
		Roles: []string{cns.RoleAdmin},
	})

	cfg, err := app.ucs.ConfigGet(bgCtx)
	require.Nil(t, err)
	require.NotNil(t, cfg)
	require.Empty(t, cfg.Contacts.Phone)
	require.Empty(t, cfg.Contacts.Email)

	err = app.ucs.ConfigSet(bgCtx, &entities.ConfigSt{})
	require.Equal(t, errs.NotAuthorized, err)

	err = app.ucs.ConfigSet(guestCtx, &entities.ConfigSt{})
	require.Equal(t, errs.PermissionDenied, err)

	err = app.ucs.ConfigSet(admCtx, &entities.ConfigSt{
		Contacts: entities.ConfigContactsSt{
			Phone: "71230000001",
			Email: "qwe@asd.com",
		},
	})
	require.Nil(t, err)

	cfg, err = app.ucs.ConfigGet(bgCtx)
	require.Nil(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, "71230000001", cfg.Contacts.Phone)
	require.Equal(t, "qwe@asd.com", cfg.Contacts.Email)
}
