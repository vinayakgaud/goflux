package goerror

import "fmt"

type Code string

const (
	CodeInvalidStateTransition Code = "INVALID_STATE_TRANSITION"
	CodeServerStartup          Code = "SERVER_STARTUP"
	CodeServerShutdown         Code = "SERVER_SHUTDOWN"
	CodeConnectionLimit        Code = "CONNECTION_LIMIT"
	CodeInvalidMessage         Code = "INVALID_MESSAGE"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Wrap(code Code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
