package queries

import (
	"github.com/go-toschool/sicily/graph"
	"github.com/graphql-go/graphql"
)

func Queries(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Queries",
		Fields: graphql.Fields{
			"getSession": GetSession(ctx),
			"getTalk":    GetTalk(ctx),
			"getTalks":   GetTalks(ctx),
			"getUser":    GetUser(ctx),
			"getUsers":   GetUsers(ctx),
			"me":         Me(ctx),
		},
	})
}
