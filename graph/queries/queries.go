package queries

import (
	"github.com/go-toschool/sicily/graph"
	"github.com/graphql-go/graphql"
)

// Queries register graph queries
func Queries(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Queries",
		Fields: graphql.Fields{
			"talk":  GetTalk(ctx),
			"talks": GetTalks(ctx),
			"user":  GetUser(ctx),
			"users": GetUsers(ctx),
		},
	})
}
