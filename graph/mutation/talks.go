package mutation

import (
	"context"
	"errors"

	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/graphql-go/graphql"
)

// CreateTalk create a talk in remote service.
func CreateTalk(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.Talk,
		Description: "Create talk",
		Args: graphql.FieldConfigArgument{
			"description": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			email, ok := params.Args["email"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			fullname, ok := params.Args["fullname"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.CreateRequest{
				Data: &citizens.Citizen{
					Email:    email,
					Fullname: fullname,
				},
			}

			u, err := ctx.UserService.Create(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return u.Data, nil
		},
	}
}

// RegisterTalk register a user into talk.
func RegisterTalk(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.Talk,
		Description: "Register a user into talk.",
		Args: graphql.FieldConfigArgument{
			"talk_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, ok := params.Args["talk_id"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			_, ok = params.Args["user_id"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.CreateRequest{
				// Data: &helenia.Citizen{
				// 	TalkId: talkID,
				// 	UserId: userID,
				// },
			}

			u, err := ctx.UserService.Create(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return u.Data, nil
		},
	}
}

// Talks expose UserQuery
func Talks(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "TalkMutations",
		Fields: graphql.Fields{
			"createTalk":   CreateTalk(ctx),
			"registerTalk": RegisterTalk(ctx),
		},
	})
}
