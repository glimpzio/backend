package graph

import (
	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/connections"
	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger            *misc.Logger
	ProfileService    *profile.ProfileService
	ConnectionService *connections.ConnectionService
	Auth0Config       *auth.Auth0Config
	Domain            string
}
