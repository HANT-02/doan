package x_error

import "fmt"

type Error struct {
	errCode string
}

func NewError(errCode string) *Error {
	return &Error{
		errCode: errCode,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.errCode)
}

func (e *Error) ErrCode() string {
	return e.errCode
}
