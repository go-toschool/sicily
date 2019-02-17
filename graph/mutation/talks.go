package mutation

import (
	"context"
	"errors"
	"time"

	"github.com/go-toschool/platon/talks"

	"github.com/go-toschool/helenia/assistants"
	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/graphql-go/graphql"
)

// CreateTalk create a talk in remote service.
func CreateTalk(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.Talk,
		Description: "Create talk",
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"repository": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"date": &graphql.ArgumentConfig{
				Type: graphql.DateTime,
			},
			"tags": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"user_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			title, ok := params.Args["title"].(string)
			if !ok {
				return nil, errors.New("Invalid title")
			}

			description, ok := params.Args["description"].(string)
			if !ok {
				return nil, errors.New("Invalid description")
			}

			repository, ok := params.Args["repository"].(string)
			if !ok {
				return nil, errors.New("Invalid repository")
			}

			date, ok := params.Args["date"].(time.Time)
			if !ok {
				return nil, errors.New("Invalid date")
			}

			tags, ok := params.Args["tags"].(string)
			if !ok {
				return nil, errors.New("Invalid tags")
			}

			userID, ok := params.Args["user_id"].(string)
			if !ok {
				return nil, errors.New("Invalid user_id")
			}

			ctxb := context.Background()
			opts := &talks.CreateRequest{
				Talk: &talks.Talk{
					Title:       title,
					Description: description,
					Repository:  repository,
					Date:        date.Unix(),
					Tags:        tags,
					UserId:      userID,
				},
			}

			t, err := ctx.TalkService.Create(ctxb, opts)
			if err != nil {
				return nil, err
			}

			return t.Talk, nil
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
			"user_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			talkID, ok := params.Args["talk_id"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			userID, ok := params.Args["user_id"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxa := context.Background()
			optsa := &talks.GetRequest{
				TalkId: talkID,
			}
			t, err := ctx.TalkService.Get(ctxa, optsa)
			if err != nil {
				return nil, err
			}

			opts := &assistants.CreateRequest{
				Data: &assistants.Assistant{
					Speaker:   t.GetTalk().GetUserId(),
					Assistant: userID,
					TalkId:    talkID,
				},
			}

			u, err := ctx.AssistantsService.Create(ctxa, opts)
			if err != nil {
				return nil, err
			}

			return u.Data, nil
		},
	}
}
