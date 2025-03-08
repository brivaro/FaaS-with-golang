package functions

import "errors"

var (
	ErrFunctionNotFound = errors.New("function not found")
	ErrFunctionNotPulled = errors.New("the image specified for the function can't be pulled from DockerHub")
	ErrUnauthorized     = errors.New("unauthorized to access this function")
	ErrInvalidRequest   = errors.New("invalid function request")
)
