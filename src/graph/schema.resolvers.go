package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"math"
	"net/http"

	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/graph/model"
	"github.com/glimpzio/backend/misc"
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
func (r *mutationResolver) ConnectByEmail(ctx context.Context, inviteID string, email string, subscribe bool) (*model.EmailConnection, error) {
	emailConnection, err := r.ConnectionService.ConnectByEmail(inviteID, email, subscribe)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrOperationFailed
	}

	r.Logger.InfoLog.Printf("connected by email for invite %s", inviteID)

	return &model.EmailConnection{
		ID:          emailConnection.Id,
		UserID:      emailConnection.UserId,
		Email:       emailConnection.Email,
		ConnectedAt: int(emailConnection.ConnectedAt.Unix()),
	}, nil
}

// Authenticated is the resolver for the authenticated field.
func (r *queryResolver) Authenticated(ctx context.Context) (bool, error) {
	middleware := auth.GetMiddleware(ctx)

	return middleware.Token != nil, nil
}

// Authenticate is the resolver for the authenticate field.
func (r *queryResolver) Authenticate(ctx context.Context, code string) (bool, error) {
	gc, err := misc.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	tkn, err := auth.ExchangeAuthCode(r.Auth0Config, code)
	if err != nil {
		return false, err
	}

	gc.SetSameSite(http.SameSiteStrictMode)
	gc.SetCookie(auth.ACCESS_TOKEN_COOKIE, tkn.AccessToken, tkn.ExpiresIn, "/", r.Domain, true, true)
	gc.SetCookie(auth.REFRESH_TOKEN_COOKIE, tkn.RefreshToken, math.MaxInt, "/", r.Domain, true, true)

	return true, nil
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

// EmailConnections is the resolver for the emailConnections field.
func (r *queryResolver) EmailConnections(ctx context.Context) ([]*model.EmailConnection, error) {
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

	rawEmailConnections, err := r.ConnectionService.GetEmailConnections(user.Id)
	if err != nil {
		r.Logger.ErrorLog.Println(err)

		return nil, ErrGetResourceFailed
	}

	r.Logger.InfoLog.Printf("retrieved email connections for user %s", user.Id)

	emailConnections := []*model.EmailConnection{}
	for _, rawEmailConnection := range rawEmailConnections {
		emailConnections = append(emailConnections, &model.EmailConnection{
			ID:          rawEmailConnection.Id,
			UserID:      rawEmailConnection.UserId,
			Email:       rawEmailConnection.Email,
			ConnectedAt: int(rawEmailConnection.ConnectedAt.Unix()),
		})
	}

	return emailConnections, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
