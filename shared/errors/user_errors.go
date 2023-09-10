package shared_errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid user credentials")

	ErrInvalidPayload = errors.New("invalid payload")

	ErrInvalidEmail = errors.New("please enter valid email")

	ErrInvalidContact = errors.New("please enter valid contact")

	ErrWeakPassword = errors.New("please enter strong password")
)
