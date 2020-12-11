package errs

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ServiceNA        = Err("server_not_available")
	NotAuthorized    = Err("not_authorized")
	PermissionDenied = Err("permission_denied")
	ObjectNotFound   = Err("object_not_found")
)
