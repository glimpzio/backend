package profile

import "errors"

var ErrInviteExpired = errors.New("the invite has expired")
var ErrInvalidUser = errors.New("the user does not exist")
var ErrInvalidInvite = errors.New("the invite does not exist")
