package profile

import "errors"

var ErrInviteExpired = errors.New("the invite has expired")
var ErrDoesNotExist = errors.New("the resource does not exist")
