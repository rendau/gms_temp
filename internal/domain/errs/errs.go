package errs

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ContextCancelled = Err("context_cancelled")
	ServiceNA        = Err("server_not_available")
)
