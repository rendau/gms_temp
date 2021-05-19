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
	usrName := "Name"
	usrTypeId := cns.UsrTypeUndefined

	err := app.ucs.ProfileSendPhoneValidatingCode(
		bgCtx,
		usrPhone,
		true,
	)
	errIsEqual(t, err, errs.PhoneNotExists)

	usrId, err := app.core.Usr.Create(bgCtx, &entities.UsrCUSt{
		TypeId: util.NewInt(usrTypeId),
		Phone:  &usrPhone,
		Name:   &usrName,
	})
	require.Nil(t, err)
	require.Greater(t, usrId, int64(0))

	_, _, err = app.ucs.ProfileAuth(
		bgCtx,
		&entities.PhoneAndSmsCodeSt{
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
		&entities.PhoneAndSmsCodeSt{
			Phone:   usrPhone,
			SmsCode: 1234,
		},
	)
	require.NotNil(t, err)
	errIsEqual(t, err, errs.WrongSmsCode)

	id, token, err := app.ucs.ProfileAuth(
		bgCtx,
		&entities.PhoneAndSmsCodeSt{
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
	require.Equal(t, usrTypeId, ses.TypeId)

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
		&entities.PhoneAndSmsCodeSt{
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
	require.Equal(t, usrTypeId, ses.TypeId)

	usrCtx = app.ucs.ContextWithSession(context.Background(), ses)

	profile, err := app.ucs.ProfileGet(usrCtx)
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, usrId, profile.Id)
	require.Equal(t, usrPhone, profile.Phone)
	require.Equal(t, usrName, profile.Name)
	require.Equal(t, usrTypeId, profile.TypeId)
}

func TestReg(t *testing.T) {
	prepareDbForNewTest()

	var smsCode int

	cases := []struct {
		phone       string
		name        string
		sendSmsCode bool
		smsE        error
		e           error
	}{
		{
			phone:       "73330000001",
			name:        "tstName",
			sendSmsCode: true,
		},
		{
			phone:       "733300001",
			name:        "tstName",
			sendSmsCode: false,
			e:           errs.BadPhoneFormat,
		},
		{
			phone:       "73330000001",
			name:        "tstName",
			sendSmsCode: true,
			e:           errs.PhoneExists,
		},
	}

	for cI, c := range cases {
		smsCode = 0

		if c.sendSmsCode {
			err := app.ucs.ProfileSendPhoneValidatingCode(
				context.Background(),
				c.phone,
				false,
			)
			errIsEqual(t, err, c.smsE, cI)

			if err == nil {
				smsCode = app.sms.PullCode()
				require.Greater(t, smsCode, 0, cI)
			}
		}

		id, _, err := app.ucs.ProfileReg(
			context.Background(),
			&entities.UsrRegReqSt{
				PhoneAndSmsCodeSt: entities.PhoneAndSmsCodeSt{
					Phone:   c.phone,
					SmsCode: smsCode,
				},
				Name: &c.name,
			},
		)
		errIsEqual(t, err, c.e, cI)
		if c.e == nil {
			require.Greater(t, id, int64(0))
			require.Nil(t, err, cI)
		}
	}

	err := app.ucs.ProfileSendPhoneValidatingCode(
		context.Background(),
		"73330000045",
		false,
	)
	require.Nil(t, err)

	smsCode = app.sms.PullCode()
	require.Greater(t, smsCode, 0)

	id, token, err := app.ucs.ProfileReg(
		context.Background(),
		&entities.UsrRegReqSt{
			PhoneAndSmsCodeSt: entities.PhoneAndSmsCodeSt{
				Phone:   "73330000045",
				SmsCode: smsCode,
			},
			Name: util.NewString("Name"),
		},
	)
	require.Nil(t, err)

	usrCtx := app.ucs.ContextWithSession(context.Background(), app.ucs.SessionGet(context.Background(), token))

	profile, err := app.ucs.ProfileGet(usrCtx)
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, id, profile.Id)
	require.Equal(t, "73330000045", profile.Phone)
	require.Equal(t, cns.UsrTypeUndefined, profile.TypeId)
	require.Equal(t, "Name", profile.Name)
}

func TestGetProfile(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()

	_, err := app.ucs.ProfileGet(bgCtx)
	errIsEqual(t, err, errs.NotAuthorized)

	profile, err := app.ucs.ProfileGet(ctxWithSes(t, nil, admId))
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, admId, profile.Id)
	require.Equal(t, cns.UsrTypeAdmin, profile.TypeId)
	require.Equal(t, admPhone, profile.Phone)
	require.Equal(t, admName, profile.Name)
}

func TestChangePhone(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()
	admCtx := ctxWithSes(t, nil, admId)

	usrId, err := app.ucs.UsrCreate(admCtx, &entities.UsrCUSt{
		TypeId: util.NewInt(cns.UsrTypeUndefined),
		Phone:  util.NewString("72340000001"),
		Name:   util.NewString("Name"),
	})
	require.Nil(t, err)

	usrCtx := ctxWithSes(t, nil, usrId)

	profile, err := app.ucs.ProfileGet(usrCtx)
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, "72340000001", profile.Phone)

	err = app.ucs.ProfileChangePhone(usrCtx, &entities.PhoneAndSmsCodeSt{
		Phone:   "72340000002",
		SmsCode: 1234,
	})
	errIsEqual(t, err, errs.SmsHasNotSentToPhone)

	err = app.ucs.ProfileSendPhoneValidatingCode(bgCtx, "72340000002", false)
	require.Nil(t, err)

	smsCode := app.sms.PullCode()

	err = app.ucs.ProfileChangePhone(usrCtx, &entities.PhoneAndSmsCodeSt{
		Phone:   "72340000002",
		SmsCode: 1234,
	})
	errIsEqual(t, err, errs.WrongSmsCode)

	err = app.ucs.ProfileChangePhone(usrCtx, &entities.PhoneAndSmsCodeSt{
		Phone:   "72340000002",
		SmsCode: smsCode,
	})
	require.Nil(t, err)

	profile, err = app.ucs.ProfileGet(usrCtx)
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, "72340000002", profile.Phone)

	err = app.ucs.ProfileSendPhoneValidatingCode(bgCtx, admPhone, false)
	require.Nil(t, err)

	smsCode = app.sms.PullCode()

	err = app.ucs.ProfileChangePhone(usrCtx, &entities.PhoneAndSmsCodeSt{
		Phone:   admPhone,
		SmsCode: smsCode,
	})
	errIsEqual(t, err, errs.PhoneExists)

	err = app.ucs.ProfileSendPhoneValidatingCode(bgCtx, "72340000002", false)
	require.Nil(t, err)

	smsCode = app.sms.PullCode()

	err = app.ucs.ProfileChangePhone(usrCtx, &entities.PhoneAndSmsCodeSt{
		Phone:   "72340000002",
		SmsCode: smsCode,
	})
	require.Nil(t, err)

	profile, err = app.ucs.ProfileGet(usrCtx)
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, "72340000002", profile.Phone)
}

func TestSessionChange(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()
	admCtx := ctxWithSes(t, nil, admId)

	usrId, err := app.ucs.UsrCreate(admCtx, &entities.UsrCUSt{
		TypeId: util.NewInt(cns.UsrTypeUndefined),
		Phone:  util.NewString("72340000001"),
		Name:   util.NewString("Name"),
	})
	require.Nil(t, err)

	usrToken, err := app.core.Usr.GetOrCreateToken(bgCtx, usrId)
	require.Nil(t, err)

	ses := app.ucs.SessionGet(bgCtx, usrToken)
	require.NotNil(t, ses)
	require.Equal(t, usrId, ses.ID)
	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)

	err = app.ucs.UsrUpdate(admCtx, usrId, &entities.UsrCUSt{
		TypeId: util.NewInt(cns.UsrTypeAdmin),
	})
	require.Nil(t, err)

	ses = app.ucs.SessionGet(bgCtx, usrToken)
	require.NotNil(t, ses)
	require.Equal(t, usrId, ses.ID)
	require.Equal(t, cns.UsrTypeAdmin, ses.TypeId)
}

func TestSessionDelete(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()
	admCtx := ctxWithSes(t, nil, admId)

	usrId, err := app.ucs.UsrCreate(admCtx, &entities.UsrCUSt{
		TypeId: util.NewInt(cns.UsrTypeUndefined),
		Phone:  util.NewString("72340000001"),
		Name:   util.NewString("Name"),
	})
	require.Nil(t, err)

	usrToken, err := app.core.Usr.GetOrCreateToken(bgCtx, usrId)
	require.Nil(t, err)

	ses := app.ucs.SessionGet(bgCtx, usrToken)
	require.NotNil(t, ses)
	require.Equal(t, usrId, ses.ID)
	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)

	err = app.ucs.ProfileLogout(ctxWithSes(t, nil, usrId))
	require.Nil(t, err)

	ses = app.ucs.SessionGet(bgCtx, usrToken)
	require.Equal(t, int64(0), ses.ID)

	usrToken, err = app.core.Usr.GetOrCreateToken(bgCtx, usrId)
	require.Nil(t, err)

	ses = app.ucs.SessionGet(bgCtx, usrToken)
	require.NotNil(t, ses)
	require.Equal(t, usrId, ses.ID)
	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)

	err = app.ucs.UsrDelete(admCtx, usrId)
	require.Nil(t, err)

	ses = app.ucs.SessionGet(bgCtx, usrToken)
	require.Equal(t, int64(0), ses.ID)
}