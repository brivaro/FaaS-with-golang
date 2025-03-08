package executor

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized access to function")
	ErrInvalidRequest = errors.New("invalid execution request")
	ErrMaxMessages    = errors.New("maximum number of execution request exceeded")
)
