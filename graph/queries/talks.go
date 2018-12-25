package queries

import (
	"context"
	"errors"

	"github.com/go-toschool/platon/talks"
	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/graphql-go/graphql"
)

// GetTalk fill graphql Field with data from plato service.
func GetTalk(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.Talk,
		Description: "Get talk by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "return taks information by id",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			id, ok := params.Args["id"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &talks.GetRequest{
				TalkId: id,
			}
			u, err := ctx.TalkService.Get(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return u.Talk, nil
		},
	}
}

// GetTalks get a collection of talks
func GetTalks(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types.Talk),
		Description: "Get collection of talks",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctxb := context.Background()
			opts := &talks.SelectRequest{}

			uu, err := ctx.TalkService.Select(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return uu.Talk, nil
		},
	}
}

// Talks expose talks service
func Talks(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "TalksQueries",
		Fields: graphql.Fields{
			"getTalk":  GetTalk(ctx),
			"getTalks": GetTalks(ctx),
		},
	})
}
