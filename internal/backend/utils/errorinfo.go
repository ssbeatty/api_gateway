package utils

import (
	"fmt"
)

const (
	StandardFormat = "2006-01-02 15:04:05"
	PreciseFormat  = "2006-01-02 15:04:05.000"
)

type AppError struct {
	Code    int32
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) OutputCode() int32 {
	return e.Code
}

func (e AppError) AppendData(data string) *AppError {
	e.Message = fmt.Sprintf(e.Message+"(%s)", data)
	return &e
}

func NewError(code int32, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

var (
	ErrSuccess        = &AppError{200, "Operate success"}
	ErrNotFound       = &AppError{404, "Not found"}
	ErrSysException   = &AppError{100001, "System error"}
	ErrBadRequest     = &AppError{100002, "Request error"}
	ErrDBError        = &AppError{100003, "Database error"}
	ErrAuthToken      = &AppError{100004, "Token expired"}
	ErrArgument       = &AppError{100005, "Argument error"}
	ErrNoPermission   = &AppError{100006, "No permission"}
	ErrGenerateToken  = &AppError{100007, "Generate token fail"}
	ErrNameOrPassword = &AppError{100008, "User name or password error"}
	ErrPasswordEmpty  = &AppError{100009, "Password cannot be empty"}
	ErrUserNameExists = &AppError{100010, "User name already exists"}
	ErrUserNotExist   = &AppError{100011, "User not exist"}
)
