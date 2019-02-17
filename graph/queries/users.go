package queries

import (
	"context"
	"errors"

	"github.com/go-toschool/platon/talks"
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
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			userID, ok := params.Args["id"].(string)
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
		Description: "Full user data",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
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

			topts := &talks.GetRequest{
				UserId: u.GetData().GetId(),
			}
			t, err := ctx.TalkService.Get(ctxb, topts)
			if err != nil {
				return nil, err
			}

			data := new(struct {
				User  *citizens.Citizen
				Talks []*talks.Talk
			})

			talks := make([]*talks.Talk, 0)
			talks = append(talks, t.GetTalk())
			data.User = u.GetData()
			data.Talks = talks

			return data, nil
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
