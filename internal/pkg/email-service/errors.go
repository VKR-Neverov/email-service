package emailservice

import "errors"

var (
	ErrCodeNotFound  = errors.New("code not found")
	ErrCodeIsInvalid = errors.New("invalid code")
)
