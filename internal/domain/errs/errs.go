package errs

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	BadJson           = Err("bad_json")
	ServiceNA         = Err("server_not_available")
	NotAuthorized     = Err("not_authorized")
	PermissionDenied  = Err("permission_denied")
	ObjectNotFound    = Err("object_not_found")
	IncorrectPageSize = Err("incorrect_page_size")

	TypeRequired         = Err("type_required")
	BadType              = Err("bad_type")
	PhoneRequired        = Err("phone_required")
	BadPhoneFormat       = Err("bad_phone_format")
	PhoneNotExists       = Err("phone_not_exists")
	PhoneExists          = Err("phone_exists")
	SmsSendLimitReached  = Err("sms_send_limit_reached")
	SmsSendTooFrequent   = Err("sms_send_too_frequent")
	SmsSendFail          = Err("sms_send_fail")
	SmsHasNotSentToPhone = Err("sms_has_not_sent_to_phone")
	WrongSmsCode         = Err("wrong_sms_code")
	NameRequired         = Err("name_required")
	BadName              = Err("bad_name")
)
