package errors

import "errors"

type DomainError struct {
	ErrType ErrorType
	err     error
}

type ErrorType int

const (
	InvalidParameter ErrorType = iota
	UnPemitedOperation
	AlreadyExist
	RepositoryError
	QueryError
	QueryDataNotFoundError
	ErrorUnknown
)

func (e *DomainError) Error() string {
	if e == nil {
		return ""
	}
	return e.err.Error()
}

func (e *DomainError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *DomainError) GetType() ErrorType {
	if e == nil {
		return ErrorUnknown
	}
	return e.ErrType
}

func NewDomainError(errType ErrorType, message string) *DomainError {
	return &DomainError{ErrType: errType, err: errors.New(message)}
}
