package app

import (
	"errors"
	"fmt"
)

type ErrCode string

const (
	ErrCodeInvalid  ErrCode = "invalid"
	ErrCodeNotFound ErrCode = "not_found"
	ErrCodeConflict ErrCode = "conflict"
	ErrCodeInternal ErrCode = "internal"
)

type Error struct {
	code    ErrCode
	message string
}

func NewError(e ErrCode, m string) *Error {
	return &Error{
		code:    e,
		message: m,
	}
}

func NewErrorf(e ErrCode, format string, args ...any) *Error {
	return &Error{
		code:    e,
		message: fmt.Sprintf(format, args...),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("internal error: code: %s, message: %s", e.code, e.message)
}

func GetErrorMessage(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.message
	}
	return err.Error()
}

func GetErrorCode(err error) ErrCode {
	var e *Error
	if errors.As(err, &e) {
		return e.code
	}
	return ErrCodeInternal
}
