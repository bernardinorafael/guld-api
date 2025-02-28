package errors

import (
	"encoding/json"
	"fmt"
	"time"
)

type ErrorCode string

type Field struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

const (
	DuplicatedField         ErrorCode = "DUPLICATED_FIELD"
	AccessTokenUnauthorized ErrorCode = "ACCESS_TOKEN_UNAUTHORIZED"
	InternalServerError     ErrorCode = "INTERNAL_SERVER_ERROR"
	BadRequest              ErrorCode = "BAD_REQUEST"
	DBResourceNotFound      ErrorCode = "DB_RESOURCE_NOT_FOUND"
	ResourceNotFound        ErrorCode = "RESOURCE_NOT_FOUND"
	ValidationField         ErrorCode = "VALIDATION_FIELD"
	InvalidPassword         ErrorCode = "INVALID_PASSWORD"
	ResourceAlreadyTaken    ErrorCode = "RESOURCE_ALREADY_TAKEN"
	LockedResource          ErrorCode = "LOCKED_RESOURCE"
	InvalidCredentials      ErrorCode = "INVALID_CREDENTIALS"
	DisabledAccount         ErrorCode = "DISABLED_ACCOUNT"
	ExpiredLink             ErrorCode = "EXPIRED_LINK"
	InvalidDeletion         ErrorCode = "INVALID_DELETION"
	MaxLimitResourceReached ErrorCode = "MAX_LIMIT_RESOURCE_REACHED"
	EmailNotVerified        ErrorCode = "EMAIL_NOT_VERIFIED"
	PhoneNotVerified        ErrorCode = "PHONE_NOT_VERIFIED"
	Expired                 ErrorCode = "EXPIRED"
)

type ApplicationError struct {
	HTTPCode  int       `json:"-"`
	Err       error     `json:"-"`
	Code      ErrorCode `json:"code"`
	Msg       string    `json:"message"`
	Timestamp int64     `json:"timestamp"`
	Fields    []Field   `json:"fields"`
}

func newApplicationError(
	httpCode int,
	code ErrorCode,
	msg string,
	err error,
	fields []Field,
) ApplicationError {
	if fields == nil {
		fields = []Field{}
	}
	return ApplicationError{
		HTTPCode:  httpCode,
		Err:       err,
		Code:      code,
		Msg:       msg,
		Timestamp: time.Now().Unix(),
		Fields:    fields,
	}
}

// AddField adds a field to the error
func (e *ApplicationError) AddField(field string, msg string) {
	e.Fields = append(
		e.Fields,
		Field{
			Field: field,
			Msg:   msg,
		},
	)
}

// RemoveField removes a field from the error
func (e *ApplicationError) RemoveField(field string) {
	for i, f := range e.Fields {
		if f.Field == field {
			e.Fields = append(e.Fields[:i], e.Fields[i+1:]...)
			break
		}
	}
}

func (e ApplicationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Msg, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e ApplicationError) JSON() string {
	out, _ := json.MarshalIndent(e, "", "  ")
	return string(out)
}

func (e ApplicationError) Unwrap() error   { return e.Err }
func (e ApplicationError) StatusCode() int { return e.HTTPCode }
