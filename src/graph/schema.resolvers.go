package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/connections"
	"github.com/glimpzio/backend/graph/model"
	"github.com/glimpzio/backend/profile"
)

// UpsertUser is the resolver for the upsertUser field.
func (r *mutationResolver) UpsertUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	middleware := auth.GetMiddleware(ctx)
	if middleware.Token == nil {
		r.Logger.ErrorLog.Println(ErrMissingAuthHeader)

		return nil, ErrMissingAuthHeader
	}

	user, err := r.ProfileService.UpsertUser(middleware.Token.AuthId, &profile.NewUser{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		PersonalEmail:  input.Email,
		Bio:            input.Bio,
		ProfilePicture: input.ProfilePicture,
		Profile: &profile.Profile{
			Email:    input.Profile.Email,
			Phone:    input.Profile.Phone,
			Website:  input.Profile.Website,
			LinkedIn: input.Profile.Linkedin,
		},
	})
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrCreateResourceFailed
	}

	r.Logger.InfoLog.Printf("upserted user %s", user.Id)

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

// CreateInvite is the resolver for the createInvite field.
func (r *mutationResolver) CreateInvite(ctx context.Context) (*model.Invite, error) {
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

	invite, err := r.ProfileService.CreateInvite(user.Id)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrCreateResourceFailed
	}

	r.Logger.InfoLog.Printf("created invite %s", invite.Id)

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

// ConnectByEmail is the resolver for the connectByEmail field.
func (r *mutationResolver) ConnectByEmail(ctx context.Context, inviteID string, email string, subscribe bool) (*model.CustomConnection, error) {
	rawConnection, err := r.ConnectionService.ConnectByEmail(inviteID, email, subscribe)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrOperationFailed
	}

	r.Logger.InfoLog.Printf("created custom connection from email for invite %s", inviteID)

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

// UpsertCustomConnection is the resolver for the upsertCustomConnection field.
func (r *mutationResolver) UpsertCustomConnection(ctx context.Context, id *string, customConnection model.NewCustomConnection) (*model.CustomConnection, error) {
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

	rawConnection, err := r.ConnectionService.UpsertCustomConnection(user.Id, id, &connections.NewCustomConnection{
		FirstName: customConnection.FirstName,
		LastName:  customConnection.LastName,
		Notes:     customConnection.Notes,
		Email:     customConnection.Email,
		Phone:     customConnection.Phone,
		Website:   customConnection.Website,
		LinkedIn:  customConnection.Linkedin,
	})
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrOperationFailed
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

// DeleteCustomConnection is the resolver for the deleteCustomConnection field.
func (r *mutationResolver) DeleteCustomConnection(ctx context.Context, id string) (*model.CustomConnection, error) {
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

	rawConnection, err = r.ConnectionService.DeleteCustomConnection(id)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrOperationFailed
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

// UploadProfilePicture is the resolver for the uploadProfilePicture field.
func (r *mutationResolver) UploadProfilePicture(ctx context.Context, file graphql.Upload) (string, error) {
	middleware := auth.GetMiddleware(ctx)
	if middleware.Token == nil {
		r.Logger.ErrorLog.Println(ErrMissingAuthHeader)

		return "", ErrMissingAuthHeader
	}

	return r.ImageUploader.ResizeAndUploadFile(file.File, 400, 400, "PROFILE_PICTURE", middleware.Token.AuthId)
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
func (r *queryResolver) CustomConnections(ctx context.Context) ([]*model.CustomConnection, error) {
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

	rawConnections, err := r.ConnectionService.GetCustomConnections(user.Id)
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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
