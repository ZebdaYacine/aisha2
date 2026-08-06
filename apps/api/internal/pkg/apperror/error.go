package apperror

import "errors"

type FieldErrors map[string]string

type Error struct {
	Code    Code
	Message string
	Fields  FieldErrors
	cause   error
}

func New(code Code, message string, fields FieldErrors) *Error {
	return &Error{Code: code, Message: message, Fields: fields}
}

func Wrap(cause error, code Code, message string) *Error {
	return &Error{Code: code, Message: message, cause: cause}
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.cause }

func Is(err, target error) bool { return errors.Is(err, target) }
