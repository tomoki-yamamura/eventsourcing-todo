package errors

type Error struct {
	ErrCode ErrCode
	Err     error
	Message string
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}
