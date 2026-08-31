package apperror

import "fmt"

type AppError struct {
	Code     int
	UserMsg  string
	Internal error
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("code=%d msg=%s err=%v", e.Code, e.UserMsg, e.Internal)
	}
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.UserMsg)
}

func New(code int, msg string) *AppError {
	return &AppError{Code: code, UserMsg: msg}
}

func Wrap(code int, msg string, err error) *AppError {
	return &AppError{Code: code, UserMsg: msg, Internal: err}
}

func Code(err error) int {
	if e, ok := err.(*AppError); ok {
		return e.Code
	}
	return CodeDBError
}

func UserMsg(err error) string {
	if e, ok := err.(*AppError); ok {
		return e.UserMsg
	}
	return "系统错误"
}
