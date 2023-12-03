package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/graph/model"
	"github.com/google/uuid"
)

// Upload is the resolver for the upload field.
func (r *queryResolver) Upload(ctx context.Context) (*model.UploadLink, error) {
	middleware := auth.GetMiddleware(ctx)
	if middleware.Token == nil {
		r.Logger.ErrorLog.Println(ErrMissingAuthHeader)

		return nil, ErrMissingAuthHeader
	}

	key := uuid.New().String()

	uploadUrl, publicUrl, err := r.ImageUploader.GetUploadLink(key)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrOperationFailed
	}

	r.Logger.InfoLog.Printf("created upload link for file %s", key)

	return &model.UploadLink{
		UploadURL: uploadUrl,
		PublicURL: publicUrl,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	middleware := auth.GetMiddleware(ctx)
	if middleware.Token == nil {
		r.Logger.ErrorLog.Println(ErrMissingAuthHeader)

		return nil, ErrMissingAuthHeader
	}

	user, err := r.ProfileService.GetUserByAuthId(middleware.Token.AuthId)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrGetResourceFailed
	}

	r.Logger.InfoLog.Printf("retrieved data for user %s", user.Id)

	return &model.User{
		ID:             user.Id,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Bio:            user.Bio,
		ProfilePicture: user.ProfilePicture,
		Profile: &model.Profile{
			Email:    user.Profile.Email,
			Phone:    user.Profile.Phone,
			Website:  user.Profile.Website,
			Linkedin: user.Profile.LinkedIn,
		},
	}, nil
}

// Invite is the resolver for the invite field.
func (r *queryResolver) Invite(ctx context.Context, id string) (*model.Invite, error) {
	invite, user, err := r.ProfileService.GetInvite(id)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrGetResourceFailed
	}

	r.Logger.InfoLog.Printf("retrieved invite %s", invite.Id)

	return &model.Invite{
		ID:        invite.Id,
		UserID:    invite.UserId,
		ExpiresAt: int(invite.ExpiresAt.Unix()),
		PublicProfile: &model.PublicProfile{
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Bio:            user.Bio,
			ProfilePicture: user.ProfilePicture,
			Profile: &model.Profile{
				Email:    user.Profile.Email,
				Phone:    user.Profile.Phone,
				Website:  user.Profile.Website,
				Linkedin: user.Profile.LinkedIn,
			},
		},
	}, nil
}

// CustomConnection is the resolver for the customConnection field.
func (r *queryResolver) CustomConnection(ctx context.Context, id string) (*model.CustomConnection, error) {
	middleware := auth.GetMiddleware(ctx)
	if middleware.Token == nil {
		r.Logger.ErrorLog.Println(ErrMissingAuthHeader)

		return nil, ErrMissingAuthHeader
	}

	rawConnection, err := r.ConnectionService.GetCustomConnection(id)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrGetResourceFailed
	}

	user, err := r.ProfileService.GetUserById(rawConnection.UserId)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrGetResourceFailed
	}

	if user.AuthId != middleware.Token.AuthId {
		r.Logger.ErrorLog.Println(ErrNotAuthorized)

		return nil, ErrNotAuthorized
	}

	return &model.CustomConnection{
		ID:          rawConnection.Id,
		UserID:      rawConnection.UserId,
		ConnectedAt: int(rawConnection.ConnectedAt.Unix()),
		FirstName:   rawConnection.FirstName,
		LastName:    rawConnection.LastName,
		Notes:       rawConnection.Notes,
		Email:       rawConnection.Email,
		Phone:       rawConnection.Phone,
		Website:     rawConnection.Website,
		Linkedin:    rawConnection.LinkedIn,
	}, nil
}

// CustomConnections is the resolver for the customConnections field.
func (r *queryResolver) CustomConnections(ctx context.Context, limit int, offset int) ([]*model.CustomConnection, error) {
	middleware := auth.GetMiddleware(ctx)
	if middleware.Token == nil {
		r.Logger.ErrorLog.Println(ErrMissingAuthHeader)

		return nil, ErrMissingAuthHeader
	}

	user, err := r.ProfileService.GetUserByAuthId(middleware.Token.AuthId)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrGetResourceFailed
	}

	r.Logger.InfoLog.Printf("retrieved data for user %s", user.Id)

	rawConnections, err := r.ConnectionService.GetCustomConnections(user.Id, limit, offset)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrGetResourceFailed
	}

	r.Logger.InfoLog.Printf("retrieved connections for user %s", user.Id)

	connections := []*model.CustomConnection{}

	for _, rawConnection := range rawConnections {
		connections = append(connections, &model.CustomConnection{
			ID:          rawConnection.Id,
			UserID:      rawConnection.UserId,
			ConnectedAt: int(rawConnection.ConnectedAt.Unix()),
			FirstName:   rawConnection.FirstName,
			LastName:    rawConnection.LastName,
			Notes:       rawConnection.Notes,
			Email:       rawConnection.Email,
			Phone:       rawConnection.Phone,
			Website:     rawConnection.Website,
			Linkedin:    rawConnection.LinkedIn,
		})
	}

	return connections, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
