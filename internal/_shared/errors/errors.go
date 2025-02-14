package errors

import (
	"encoding/json"
	"net/http"
)

func NewHttpError(w http.ResponseWriter, err ApplicationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode())
	_ = json.NewEncoder(w).Encode(err)
}

func NewUnauthorizedError(msg string, err error) ApplicationError {
	httpCode := http.StatusUnauthorized
	return newApplicationError(httpCode, AccessTokenUnauthorized, msg, err, nil)
}

func NewInternalServerError(err error) ApplicationError {
	httpCode := http.StatusInternalServerError
	return newApplicationError(httpCode, InternalServerError, "something went wrong", err, nil)
}

func NewForbiddenError(msg string, code ErrorCode, err error) ApplicationError {
	httpCode := http.StatusForbidden
	return newApplicationError(httpCode, code, msg, err, nil)
}

func NewValidationFieldError(msg string, err error, fields []Field) ApplicationError {
	httpCode := http.StatusUnprocessableEntity
	return newApplicationError(httpCode, ValidationField, msg, err, fields)
}

func NewBadRequestError(msg string, err error) ApplicationError {
	httpCode := http.StatusBadRequest
	return newApplicationError(httpCode, BadRequest, msg, err, nil)
}

func NewNotFoundError(msg string, err error) ApplicationError {
	httpCode := http.StatusNotFound
	return newApplicationError(httpCode, DBResourceNotFound, msg, err, nil)
}

func NewConflictError(msg string, code ErrorCode, err error, fields []Field) ApplicationError {
	httpCode := http.StatusConflict
	return newApplicationError(httpCode, code, msg, err, fields)
}
