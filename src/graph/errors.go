package graph

import "errors"

var ErrMissingAuthHeader = errors.New("missing auth header")
var ErrNotAuthorized = errors.New("not authorized to access")
var ErrGetResourceFailed = errors.New("failed to retrieve resource")
var ErrCreateResourceFailed = errors.New("failed to create resource")
var ErrOperationFailed = errors.New("the operation failed")
var ErrUnknown = errors.New("an unknown error occured")
