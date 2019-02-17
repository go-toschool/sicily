package types

import (
	"github.com/graphql-go/graphql"
)

var User = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"full_name": &graphql.Field{
			Type: graphql.String,
		},
		"token": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.Int,
		},
		"updated_at": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

// UserWithTalks this store user information and it subscribed talks.
var UserWithTalks = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserWithTalks",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: User,
		},
		"talks": &graphql.Field{
			Type: graphql.NewList(Talk),
		},
	},
})
