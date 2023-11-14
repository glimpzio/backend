package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/glimpzio/backend/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	req, _ := ctx.Value("httpRequest").(*http.Request)
	fmt.Println(req.Header.Get(("Authorization")))

	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// CreateLink is the resolver for the createLink field.
func (r *mutationResolver) CreateLink(ctx context.Context, input *model.NewLink) (*model.Link, error) {
	panic(fmt.Errorf("not implemented: CreateLink - createLink"))
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, input *model.NewProfile) (*model.Profile, error) {
	panic(fmt.Errorf("not implemented: UpdateProfile - updateProfile"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return nil, errors.New("test")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
