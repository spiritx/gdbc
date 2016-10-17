package gdbc

import "fmt"

type DbError interface {
	Code() int
	Error() string
}

type DefaultDbError struct {
	ErrorCode int
	ErrorMsg  string
}

func NewDefaultDbError(errorCode int, errorMsg string) *DefaultDbError {
	if errorCode == 0 {
		return nil
	}
	return &DefaultDbError{errorCode, errorMsg}
}

func NewDefaultDbErrorf(errorCode int, format string, v ...interface{}) *DefaultDbError {
	if errorCode == 0 {
		return nil
	}
	errorMsg := fmt.Sprintf(format, v...)
	return &DefaultDbError{errorCode, errorMsg}
}

func (this *DefaultDbError) IsFailure() bool {
	if this == nil {
		return false
	}

	if this.ErrorCode == 0 {
		return false
	}
	return true
}

func (this *DefaultDbError) IsOk() bool {
	return !this.IsFailure()
}

func (this *DefaultDbError) Code() int {
	if this == nil {
		return 0
	}

	return this.ErrorCode
}

func (this *DefaultDbError) Error() string {
	if this == nil {
		return ""
	}

	return this.ErrorMsg
}
