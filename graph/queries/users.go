package queries

import (
	"context"
	"errors"

	"github.com/go-toschool/sicily"
	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/graphql-go/graphql"
)

const (
	// UserIDKey ...
	UserIDKey sicily.StringValueKey = "user_id"
)

// GetUser fill graphql Field with data from postgres service.
func GetUser(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.User,
		Description: "Get user by id",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			userID, ok := params.Context.Value(UserIDKey).(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.GetRequest{
				UserId: userID,
			}
			u, err := ctx.UserService.Get(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return u.Data, nil
		},
	}
}

// Me resolve user information and users talk from external services.
func Me(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.UserWithTalks,
		Description: "Fill user data",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			userID, ok := params.Context.Value(UserIDKey).(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.GetRequest{
				UserId: userID,
			}
			u, err := ctx.UserService.Get(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return u.Data, nil
		},
	}
}

// GetUsers get a collection of users
func GetUsers(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types.User),
		Description: "Get collection of users",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctxb := context.Background()
			opts := &citizens.SelectRequest{}
			uu, err := ctx.UserService.Select(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return uu.Data, nil
		},
	}
}

// Users expose UserQuery
func Users(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "UserQueries",
		Fields: graphql.Fields{
			"getUser":  GetUser(ctx),
			"getUsers": GetUsers(ctx),
			"me":       Me(ctx),
		},
	})
}
