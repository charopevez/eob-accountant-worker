package apperror

import (
	"encoding/json"
	"fmt"
)

var (
	//account status error
	ErrNotFound   = NewAppError("not found", "NS-000010", "")
	ErrNotActive  = NewAppError("account isn't active", "NS-000011", "Please check you email for activation link")
	ErrIsDeleted  = NewAppError("account is deleted", "NS-000012", "")
	ErrNotMatched = NewAppError("wrong password", "NS-000012", "")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func NewAppError(message, code, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(message, "NS-000003", "")
}

func BadRequestError(message string) *AppError {
	return NewAppError(message, "NS-000002", "something wrong with user data")
}

func systemError(developerMessage string) *AppError {
	return NewAppError("system error", "NS-000001", developerMessage)
}

func APIError(code, message, developerMessage string) *AppError {
	return NewAppError(message, code, developerMessage)
}
