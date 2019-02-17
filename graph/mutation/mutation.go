package mutation

import (
	"github.com/go-toschool/sicily/graph"
	"github.com/graphql-go/graphql"
)

func Mutations(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutations",
		Fields: graphql.Fields{
			"createTalk":   CreateTalk(ctx),
			"registerTalk": RegisterTalk(ctx),
			"updateUser":   UpdateUser(ctx),
		},
	})
}
