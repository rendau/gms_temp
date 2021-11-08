package core

import (
	"context"
	"crypto/sha512"
	"fmt"
	"math/rand"
	"time"

	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
	"github.com/rendau/gms_temp/internal/domain/errs"
	"github.com/rendau/gms_temp/internal/domain/util"
)

const (
	unValidateCacheKeyTmpl    = "phone_validate_%s"
	unValidateCacheTimeout    = 20 * time.Minute
	unValidateSmsSendLimit    = 3
	unValidateSmsSendInterval = 2 * time.Minute
)

var (
	smsFreePhones = map[string]int{ // phone and static code, code 0 - mean any code
		"70000000000": 7000,
	}
)

type Usr struct {
	r *St
}

type usrUNValidatingCacheSt struct {
	Codes        []int     `json:"codes"`
	LastSentTime time.Time `json:"lst"`
}

func NewUsr(r *St) *Usr {
	return &Usr{
		r: r,
	}
}

func (c *Usr) ValidatePhone(phone string) (string, error) {
	fPhone := util.NormalizePhone(phone)
	if !util.ValidatePhone(fPhone) {
		return phone, errs.BadPhoneFormat
	}

	return fPhone, nil
}

func (c *Usr) SendPhoneValidatingCode(ctx context.Context, phone string, errNE bool) error {
	var err error

	phone, err = c.ValidatePhone(phone)
	if err != nil {
		return err
	}

	id, err := c.GetIdForPhone(ctx, phone, errNE)
	if err != nil {
		return err
	}

	if id > 0 {
		// check if id is blocked
	}

	if c.r.noSmsCheck {
		return nil
	}

	// check if phone is sms free
	if _, ok := smsFreePhones[phone]; ok {
		return nil
	}

	var cacheValue usrUNValidatingCacheSt
	var cacheKey = fmt.Sprintf(unValidateCacheKeyTmpl, phone)

	_, err = c.r.cache.GetJsonObj(cacheKey, &cacheValue)
	if err != nil {
		return err
	}

	if len(cacheValue.Codes) >= unValidateSmsSendLimit {
		return errs.SmsSendLimitReached
	}

	if cacheValue.LastSentTime.Add(unValidateSmsSendInterval).After(time.Now()) {
		return errs.SmsSendTooFrequent
	}

	rand.Seed(time.Now().UTC().UnixNano())
	smsCode := 1000 + rand.Intn(9000)

	var smsText string

	smsText = fmt.Sprintf("%d - Используйте данный код для %s", smsCode, cns.AppName)

	ok := c.r.sms.Send(phone, smsText)
	if !ok {
		return errs.SmsSendFail
	}

	cacheValue.Codes = append(cacheValue.Codes, smsCode)
	cacheValue.LastSentTime = time.Now()

	err = c.r.cache.SetJsonObj(cacheKey, cacheValue, unValidateCacheTimeout)
	if err != nil {
		return err
	}

	if !c.r.testing {
		c.r.lg.Infow("Sms sent", "phone", phone, "sms_code", smsCode)
	}

	return nil
}

func (c *Usr) CheckPhoneValidatingCode(ctx context.Context, obj *entities.PhoneAndSmsCodeSt) error {
	if c.r.noSmsCheck {
		return nil
	}

	// check if phone is sms free
	if staticCode, ok := smsFreePhones[obj.Phone]; ok {
		if staticCode == 0 {
			return nil
		} else if obj.SmsCode != staticCode {
			return errs.WrongSmsCode
		}
	}

	var cacheKey = fmt.Sprintf(unValidateCacheKeyTmpl, obj.Phone)
	var cacheValue usrUNValidatingCacheSt

	ok, err := c.r.cache.GetJsonObj(cacheKey, &cacheValue)
	if err != nil {
		return err
	}
	if !ok {
		return errs.SmsHasNotSentToPhone
	}

	smsCodeFound := false
	for _, code := range cacheValue.Codes {
		if code == obj.SmsCode {
			smsCodeFound = true
			break
		}
	}

	if !smsCodeFound {
		return errs.WrongSmsCode
	}

	return nil
}

func (c *Usr) RemovePhoneValidatingCache(ctx context.Context, phone string) {
	if c.r.noSmsCheck {
		return
	}

	// check if phone is sms free
	if _, ok := smsFreePhones[phone]; ok {
		return
	}

	_ = c.r.cache.Del(fmt.Sprintf(unValidateCacheKeyTmpl, phone))
}

func (c *Usr) List(ctx context.Context, pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error) {
	return c.r.db.UsrList(ctx, pars)
}

func (c *Usr) ListOne(ctx context.Context, id int64) (*entities.UsrListSt, error) {
	rows, _, err := c.List(ctx, &entities.UsrListParsSt{
		Ids: util.NewSliceInt64(id),
	})
	if err != nil {
		return nil, err
	}

	if len(rows) < 1 {
		return nil, errs.ObjectNotFound
	}

	return rows[0], nil
}

func (c *Usr) Get(ctx context.Context, pars *entities.UsrGetParsSt, errNE bool) (*entities.UsrSt, error) {
	var err error

	var result *entities.UsrSt

	if pars.Id != nil || pars.Phone != nil {
		result, err = c.r.db.UsrGet(ctx, pars)
		if err != nil {
			return nil, err
		}
	}

	if result == nil {
		if errNE {
			return nil, errs.ObjectNotFound
		}

		return nil, nil
	}

	return result, nil
}

func (c *Usr) IdExists(ctx context.Context, id int64) (bool, error) {
	return c.r.db.UsrIdExists(ctx, id)
}

func (c *Usr) IdsExists(ctx context.Context, ids []int64) (bool, error) {
	return c.r.db.UsrIdsExists(ctx, ids)
}

func (c *Usr) PhoneExists(ctx context.Context, phone string, excludeId int64) (bool, error) {
	return c.r.db.UsrPhoneExists(ctx, phone, excludeId)
}

func (c *Usr) GetToken(ctx context.Context, id int64) (string, error) {
	return c.r.db.UsrGetToken(ctx, id)
}

func (c *Usr) GenerateAndSaveToken(ctx context.Context, id int64) (string, error) {
	tokenSrc := fmt.Sprintf("auth-lt-token %d %s", id, time.Now())
	token := fmt.Sprintf("%x", sha512.Sum512([]byte(tokenSrc)))

	err := c.SetToken(ctx, id, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *Usr) SetToken(ctx context.Context, id int64, v string) error {
	err := c.r.db.UsrSetToken(ctx, id, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Usr) ResetToken(ctx context.Context, id int64) error {
	return c.SetToken(ctx, id, "")
}

func (c *Usr) GetTypeId(ctx context.Context, id int64) (int, error) {
	return c.r.db.UsrGetTypeId(ctx, id)
}

func (c *Usr) GetPhone(ctx context.Context, id int64) (string, error) {
	return c.r.db.UsrGetPhone(ctx, id)
}

func (c *Usr) GetIdForPhone(ctx context.Context, phone string, errNE bool) (int64, error) {
	var err error

	phone, err = c.ValidatePhone(phone)
	if err != nil {
		return 0, nil
	}

	id, err := c.r.db.UsrGetIdForPhone(ctx, phone)
	if err != nil {
		return 0, err
	}

	if id == 0 && errNE {
		return 0, errs.PhoneNotExists
	}

	return id, nil
}

func (c *Usr) Auth(ctx context.Context, obj *entities.PhoneAndSmsCodeSt) (int64, string, error) {
	var err error

	obj.Phone, err = c.ValidatePhone(obj.Phone)
	if err != nil {
		return 0, "", err
	}

	err = c.CheckPhoneValidatingCode(ctx, obj)
	if err != nil {
		return 0, "", err
	}

	usr, err := c.Get(ctx, &entities.UsrGetParsSt{
		Phone: &obj.Phone,
	}, true)
	if err != nil {
		return 0, "", err
	}

	token, err := c.r.Session.CreateToken(&entities.Session{
		Id:     usr.Id,
		TypeId: usr.TypeId,
	})
	if err != nil {
		return 0, "", err
	}

	c.RemovePhoneValidatingCache(ctx, obj.Phone)

	return usr.Id, token, nil
}

func (c *Usr) Reg(ctx context.Context, data *entities.UsrRegReqSt) (int64, string, error) {
	var err error

	if data.Phone, err = c.ValidatePhone(data.Phone); err != nil {
		return 0, "", err
	}

	err = c.CheckPhoneValidatingCode(ctx, &data.PhoneAndSmsCodeSt)
	if err != nil {
		return 0, "", err
	}

	id, err := c.GetIdForPhone(ctx, data.Phone, false)
	if err != nil {
		return 0, "", err
	}
	if id > 0 {
		return 0, "", errs.PhoneExists
	}

	id, err = c.Create(ctx, &entities.UsrCUSt{
		TypeId: data.TypeId,
		Phone:  &data.Phone,
		Ava:    data.Ava,
		Name:   data.Name,
	})
	if err != nil {
		return 0, "", err
	}

	newUsr, err := c.Get(ctx, &entities.UsrGetParsSt{
		Id: &id,
	}, true)
	if err != nil {
		return 0, "", err
	}

	token, err := c.r.Session.CreateToken(&entities.Session{
		Id:     newUsr.Id,
		TypeId: newUsr.TypeId,
	})
	if err != nil {
		return 0, "", err
	}

	c.RemovePhoneValidatingCache(ctx, data.Phone)

	return id, token, nil
}

func (c *Usr) ValidateCU(ctx context.Context, obj *entities.UsrCUSt, id int64) error {
	forCreate := id == 0

	// type_id
	if forCreate && obj.TypeId == nil {
		return errs.TypeRequired
	}
	if obj.TypeId != nil {
		if !cns.UsrTypeIsValid(*obj.TypeId) {
			return errs.BadType
		}
	}

	// phone
	if forCreate && obj.Phone == nil {
		return errs.PhoneRequired
	}
	if obj.Phone != nil {
		var err error

		if *obj.Phone, err = c.ValidatePhone(*obj.Phone); err != nil {
			return err
		}

		exists, err := c.PhoneExists(ctx, *obj.Phone, id)
		if err != nil {
			return err
		}
		if exists {
			return errs.PhoneExists
		}
	}

	// name
	if forCreate && obj.Name == nil {
		return errs.NameRequired
	}
	if obj.Name != nil {
		if *obj.Name == "" {
			return errs.NameRequired
		} else if len([]rune(*obj.Name)) > 250 {
			return errs.BadName
		}
	}

	return nil
}

func (c *Usr) Create(ctx context.Context, obj *entities.UsrCUSt) (int64, error) {
	err := c.ValidateCU(ctx, obj, 0)
	if err != nil {
		return 0, err
	}

	// create
	newId, err := c.r.db.UsrCreate(ctx, obj)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (c *Usr) GetProfile(ctx context.Context, id int64) (*entities.UsrProfileSt, error) {
	usr, err := c.Get(ctx, &entities.UsrGetParsSt{Id: &id}, true)
	if err != nil {
		return nil, err
	}

	res := &entities.UsrProfileSt{
		UsrSt: *usr,
	}

	return res, nil
}

func (c *Usr) GetNumbers(ctx context.Context, id int64) (*entities.UsrNumbersSt, error) {
	result := &entities.UsrNumbersSt{}

	// result.NewMsgCount, err = c.r.GetNewMsgCount(ctx, id)
	// if err != nil {
	// 	return nil, err
	// }

	return result, nil
}

func (c *Usr) Update(ctx context.Context, id int64, obj *entities.UsrCUSt) error {
	err := c.ValidateCU(ctx, obj, id)
	if err != nil {
		return err
	}

	err = c.r.db.UsrUpdate(ctx, id, obj)
	if err != nil {
		return err
	}

	return nil
}

func (c *Usr) ChangePhone(ctx context.Context, id int64, obj *entities.PhoneAndSmsCodeSt) error {
	var err error

	obj.Phone, err = c.ValidatePhone(obj.Phone)
	if err != nil {
		return err
	}

	err = c.CheckPhoneValidatingCode(ctx, obj)
	if err != nil {
		return err
	}

	err = c.Update(ctx, id, &entities.UsrCUSt{
		Phone: &obj.Phone,
	})
	if err != nil {
		return err
	}

	c.RemovePhoneValidatingCache(ctx, obj.Phone)

	return nil
}

func (c *Usr) Logout(ctx context.Context, id int64) error {
	return nil
}

func (c *Usr) Delete(ctx context.Context, id int64) error {
	err := c.r.db.UsrDelete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
