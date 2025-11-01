package errors

import "errors"

type ErrCode string

const (
	Unknown          ErrCode = "U000"
	InvalidParameter ErrCode = "V001"
	UnpermittedOp    ErrCode = "A001"
	NotFound         ErrCode = "N001"
	RepositoryError  ErrCode = "R001"
	QueryError       ErrCode = "R002"
	OptimisticLock   ErrCode = "R003"
)

func (code ErrCode) New(message string) error {
	return &Error{
		ErrCode: code,
		Message: message,
		Err:     errors.New(message),
	}
}

func (code ErrCode) Wrap(err error, message string) error {
	return &Error{
		ErrCode: code,
		Message: message,
		Err:     err,
	}
}

func IsCode(err error, code ErrCode) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.ErrCode == code
	}
	return false
}
