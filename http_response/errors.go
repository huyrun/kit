package http_response

type ErrUnauthorized struct {
	msg string
}

func NewErrUnauthorized(msg string) *ErrUnauthorized {
	return & ErrUnauthorized{
		msg: msg,
	}
}

func (e *ErrUnauthorized) Error() string {
	return e.msg
}

type ErrInternal struct {
	msg string
}

func NewErrInternal(msg string) *ErrUnauthorized {
	return & ErrUnauthorized{
		msg: msg,
	}
}

func (e *ErrInternal) Error() string {
	return e.msg
}
