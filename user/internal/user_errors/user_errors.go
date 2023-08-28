package user_errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid user credentials")

	ErrInvalidPayload = errors.New("invalid payload")
)
