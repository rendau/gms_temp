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
	require.Equal(t, errs.PhoneNotExists, err)

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
	require.Equal(t, errs.SmsHasNotSentToPhone, err)

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
	require.Equal(t, errs.WrongSmsCode, err)

	accessToken, refreshToken, err := app.ucs.ProfileAuth(
		bgCtx,
		&entities.PhoneAndSmsCodeSt{
			Phone:   usrPhone,
			SmsCode: smsCode,
		},
	)
	require.Nil(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	ses := app.ucs.SessionGetFromToken(accessToken)
	require.Nil(t, err)
	require.NotNil(t, ses)
	require.Equal(t, usrId, ses.Id)
	require.Equal(t, usrTypeId, ses.TypeId)

	usrCtx := app.ucs.SessionSetToContextByToken(context.Background(), accessToken)

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
			require.Equal(t, c.smsE, err, cI)

			if err == nil {
				smsCode = app.sms.PullCode()
				require.Greater(t, smsCode, 0, cI)
			}
		}

		accessToken, _, err := app.ucs.ProfileReg(
			context.Background(),
			&entities.UsrRegReqSt{
				PhoneAndSmsCodeSt: entities.PhoneAndSmsCodeSt{
					Phone:   c.phone,
					SmsCode: smsCode,
				},
				Name: &c.name,
			},
		)
		require.Equal(t, c.e, err, cI)
		if c.e == nil {
			require.NotEmpty(t, accessToken)
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

	accessToken, _, err := app.ucs.ProfileReg(
		context.Background(),
		&entities.UsrRegReqSt{
			PhoneAndSmsCodeSt: entities.PhoneAndSmsCodeSt{
				Phone:   "73330000045",
				SmsCode: smsCode,
			},
			Ava:  util.NewString("/path_to_ava"),
			Name: util.NewString("Name"),
		},
	)
	require.Nil(t, err)

	usrCtx := app.ucs.SessionSetToContextByToken(context.Background(), accessToken)

	profile, err := app.ucs.ProfileGet(usrCtx)
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, "73330000045", profile.Phone)
	require.Equal(t, cns.UsrTypeUndefined, profile.TypeId)
	require.Equal(t, "/path_to_ava", profile.Ava)
	require.Equal(t, "Name", profile.Name)
}

func TestProfileGet(t *testing.T) {
	prepareDbForNewTest()

	bgCtx := context.Background()

	_, err := app.ucs.ProfileGet(bgCtx)
	require.Equal(t, errs.NotAuthorized, err)

	profile, err := app.ucs.ProfileGet(ctxWithSes(t, nil, admId))
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, admId, profile.Id)
	require.Equal(t, cns.UsrTypeAdmin, profile.TypeId)
	require.Equal(t, admPhone, profile.Phone)
	require.Equal(t, admName, profile.Name)
}

func TestPhoneChange(t *testing.T) {
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
	require.Equal(t, errs.SmsHasNotSentToPhone, err)

	err = app.ucs.ProfileSendPhoneValidatingCode(bgCtx, "72340000002", false)
	require.Nil(t, err)

	smsCode := app.sms.PullCode()

	err = app.ucs.ProfileChangePhone(usrCtx, &entities.PhoneAndSmsCodeSt{
		Phone:   "72340000002",
		SmsCode: 1234,
	})
	require.Equal(t, errs.WrongSmsCode, err)

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
	require.Equal(t, errs.PhoneExists, err)

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

// func TestSessionChange(t *testing.T) {
// 	prepareDbForNewTest()
//
// 	bgCtx := context.Background()
// 	admCtx := ctxWithSes(t, nil, admId)
//
// 	usrId, err := app.ucs.UsrCreate(admCtx, &entities.UsrCUSt{
// 		TypeId: util.NewInt(cns.UsrTypeUndefined),
// 		Phone:  util.NewString("72340000001"),
// 		Name:   util.NewString("Name"),
// 	})
// 	require.Nil(t, err)
//
// 	usrToken, err := app.core.Usr.GetOrCreateToken(bgCtx, usrId)
// 	require.Nil(t, err)
//
// 	ses := app.ucs.SessionGet(bgCtx, usrToken)
// 	require.NotNil(t, ses)
// 	require.Equal(t, usrId, ses.Id)
// 	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)
//
// 	err = app.ucs.UsrUpdate(admCtx, usrId, &entities.UsrCUSt{
// 		TypeId: util.NewInt(cns.UsrTypeAdmin),
// 	})
// 	require.Nil(t, err)
//
// 	ses = app.ucs.SessionGet(bgCtx, usrToken)
// 	require.NotNil(t, ses)
// 	require.Equal(t, usrId, ses.Id)
// 	require.Equal(t, cns.UsrTypeAdmin, ses.TypeId)
// }
//
// func TestSessionDelete(t *testing.T) {
// 	prepareDbForNewTest()
//
// 	bgCtx := context.Background()
// 	admCtx := ctxWithSes(t, nil, admId)
//
// 	usrId, err := app.ucs.UsrCreate(admCtx, &entities.UsrCUSt{
// 		TypeId: util.NewInt(cns.UsrTypeUndefined),
// 		Phone:  util.NewString("72340000001"),
// 		Name:   util.NewString("Name"),
// 	})
// 	require.Nil(t, err)
//
// 	usrToken, err := app.core.Usr.GetOrCreateToken(bgCtx, usrId)
// 	require.Nil(t, err)
//
// 	ses := app.ucs.SessionGet(bgCtx, usrToken)
// 	require.NotNil(t, ses)
// 	require.Equal(t, usrId, ses.Id)
// 	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)
//
// 	err = app.ucs.ProfileLogout(ctxWithSes(t, nil, usrId))
// 	require.Nil(t, err)
//
// 	ses = app.ucs.SessionGet(bgCtx, usrToken)
// 	require.Equal(t, int64(0), ses.Id)
//
// 	usrToken, err = app.core.Usr.GetOrCreateToken(bgCtx, usrId)
// 	require.Nil(t, err)
//
// 	ses = app.ucs.SessionGet(bgCtx, usrToken)
// 	require.NotNil(t, ses)
// 	require.Equal(t, usrId, ses.Id)
// 	require.Equal(t, cns.UsrTypeUndefined, ses.TypeId)
//
// 	err = app.ucs.UsrDelete(admCtx, usrId)
// 	require.Nil(t, err)
//
// 	ses = app.ucs.SessionGet(bgCtx, usrToken)
// 	require.Equal(t, int64(0), ses.Id)
// }
