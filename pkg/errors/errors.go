package errors

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExits = errors.New("user already exists")
