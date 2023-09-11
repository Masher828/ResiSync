package shared_errors

func getInternalErrorMap() func(err error) bool {

	errMap := map[error]bool{
		ErrInvalidCredentials: true,
		ErrInvalidPayload:     true,
		ErrInvalidContact:     true,
		ErrInvalidEmail:       true,
		ErrWeakPassword:       true,
		ErrInvalidOTP:         true,
	}
	return func(err error) bool {
		return errMap[err]
	}
}

var IsInternalError = getInternalErrorMap()
