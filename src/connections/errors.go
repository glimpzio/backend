package connections

import "errors"

var ErrDoesNotExist = errors.New("the resource does not exist")
var ErrNotAuthorized = errors.New("not authorized to access this resource")
