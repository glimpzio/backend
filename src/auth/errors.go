package auth

import "errors"

var ErrMissingAuthHeader = errors.New("missing auth header")
var ErrNotAuthorized = errors.New("not authorized to access")
var ErrExpired = errors.New("resource has expired")
