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
	ImageUploader     *misc.ImageUploader
}

func NewResolver(logger *misc.Logger, profileService *profile.ProfileService, connectionService *connections.ConnectionService, auth0Config *auth.Auth0Config, imageUploader *misc.ImageUploader) *Resolver {
	return &Resolver{
		Logger:            logger,
		ProfileService:    profileService,
		ConnectionService: connectionService,
		Auth0Config:       auth0Config,
		ImageUploader:     imageUploader,
	}
}
