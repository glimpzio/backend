package graph

import (
	"log"

	"github.com/glimpzio/backend/profile"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger         *log.Logger
	ProfileService *profile.ProfileService
}
