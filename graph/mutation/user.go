package mutation

import (
	"context"
	"errors"

	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/graphql-go/graphql"
)

// UpdateUser updates basic information
func UpdateUser(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.User,
		Description: "Update user by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"full_name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			id, ok := params.Args["id"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			fullName, ok := params.Args["fullName"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.UpdateRequest{
				UserId: id,
				Data: &citizens.Citizen{
					FullName: fullName,
				},
			}

			u, err := ctx.UserService.Update(ctxb, opts)
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
		Name: "UserMutations",
		Fields: graphql.Fields{
			"updateUser": UpdateUser(ctx),
		},
	})
}
