package queries

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/go-toschool/syracuse/citizens"
)

// GetUser fill graphql Field with data from postgres service.
func GetUser(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.User,
		Description: "Get user by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "return user information by id",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			id, ok := params.Args["id"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.GetRequest{
				UserId: id,
			}
			u, err := ctx.UserService.Get(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return u.Data, nil
		},
	}
}

// Users expose UserQuery
func Users(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "UserQueries",
		Fields: graphql.Fields{
			"users": GetUser(ctx),
		},
	})
}
