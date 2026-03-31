package services

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidRefreshSession = errors.New("invalid refresh session")

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}
