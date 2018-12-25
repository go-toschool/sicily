package queries

import (
	"context"
	"errors"

	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/graphql-go/graphql"
)

// GetSession fill graphql Field with data from postgres service.
func GetSession(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.User,
		Description: "Get session data by auth token",
		Args: graphql.FieldConfigArgument{
			"token": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "token to send to rpc service",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			token, ok := params.Args["token"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.GetRequest{
				UserId: token,
			}
			u, err := ctx.UserService.Get(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return u.Data, nil
		},
	}
}

// Session expose SessionQueries
func Session(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "SessionQueries",
		Fields: graphql.Fields{
			"getSession": GetSession(ctx),
		},
	})
}
