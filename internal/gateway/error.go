package gateway

type ErrCode int

const (
	ErrCodeBadInput ErrCode = iota
	ErrCodeNotFound
)

type Error struct {
	code ErrCode
	msg  string
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Code() ErrCode {
	return e.code
}

func NewError(code ErrCode, msg string) *Error {
	return &Error{code, msg}
}
