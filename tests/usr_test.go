package tests

import (
	"context"
	"testing"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
	"github.com/rendau/gms_temp/internal/domain/util"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()

	usrPhone := "76760000001"
	usrName := "name"

	err := app.ucs.ProfileSendPhoneValidatingCode(
		bgCtx,
		usrPhone,
		true,
	)
	errIsEqual(t, err, errs.PhoneNotExists)

	usrId, err := app.core.Usr.Create(bgCtx, &entities.UsrCUSt{
		TypeId: util.NewInt(cns.UsrTypeUndefined),
		Phone:  &usrPhone,
		Name:   &usrName,
	})
	require.Nil(t, err)
	require.Greater(t, usrId, int64(0))

	_, _, err = app.ucs.ProfileAuth(
		bgCtx,
		&entities.UsrAuthReqSt{
			Phone:   usrPhone,
			SmsCode: 1234,
		},
	)
	errIsEqual(t, err, errs.SmsHasNotSentToPhone)

	err = app.ucs.ProfileSendPhoneValidatingCode(
		bgCtx,
		usrPhone,
		true,
	)
	require.Nil(t, err)

	smsCode := app.sms.PullCode()
	require.Greater(t, smsCode, 0)

	_, _, err = app.ucs.ProfileAuth(
		bgCtx,
		&entities.UsrAuthReqSt{
			Phone:   usrPhone,
			SmsCode: 1234,
		},
	)
	require.NotNil(t, err)
	errIsEqual(t, err, errs.WrongSmsCode)

	id, token, err := app.ucs.ProfileAuth(
		bgCtx,
		&entities.UsrAuthReqSt{
			Phone:   usrPhone,
			SmsCode: smsCode,
		},
	)
	require.Nil(t, err)
	require.Equal(t, usrId, id)

	ses := app.ucs.SessionGet(context.Background(), token)
	require.Nil(t, err)
	require.NotNil(t, ses)
	require.Equal(t, usrId, ses.ID)
	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)

	err = app.ucs.ProfileLogout(ctxWithSes(t, nil, usrId))
	require.Nil(t, err)

	ses = app.ucs.SessionGet(context.Background(), token)
	require.NotNil(t, ses)
	require.Equal(t, int64(0), ses.ID)

	usrCtx := app.ucs.ContextWithSession(context.Background(), app.ucs.SessionGet(context.Background(), token))

	_, err = app.ucs.ProfileGet(usrCtx)
	errIsEqual(t, err, errs.NotAuthorized)

	err = app.ucs.ProfileSendPhoneValidatingCode(
		bgCtx,
		usrPhone,
		true,
	)
	require.Nil(t, err)

	smsCode = app.sms.PullCode()
	require.Greater(t, smsCode, 0)

	id, token, err = app.ucs.ProfileAuth(
		bgCtx,
		&entities.UsrAuthReqSt{
			Phone:   usrPhone,
			SmsCode: smsCode,
		},
	)
	require.Nil(t, err)
	require.Equal(t, usrId, id)

	ses = app.ucs.SessionGet(context.Background(), token)
	require.Nil(t, err)
	require.NotNil(t, ses)
	require.Equal(t, usrId, ses.ID)
	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)
}
