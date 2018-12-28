package mutation

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-toschool/palermo/auth"
	"github.com/go-toschool/sicily"

	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/types"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/graphql-go/graphql"
)

// CreateSession create user if does not exists, and a session,
// This method call grpc function that store google login token in syracusa service.
// Also, request grpc method to create a jwt token to handle session.
func CreateSession(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.User,
		Description: "Create user session",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"token": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"full_name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			email, ok := params.Args["email"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			fullName, ok := params.Args["full_name"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}
			// google token
			token, ok := params.Args["token"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.CreateRequest{
				Data: &citizens.Citizen{
					Email:    email,
					FullName: fullName,
				},
			}
			// create user in syracusa service.
			u, err := ctx.UserService.Create(ctxb, opts)
			if err != nil {
				return nil, err
			}

			id, iat, err := idWithTime()
			if err != nil {
				return nil, err
			}
			sess := &auth.CreateRequest{
				Data: &auth.Session{
					Id:        id,
					UserId:    u.Data.Id,
					Email:     email,
					Token:     token,
					CreatedAt: iat.Unix(),
					UpdatedAt: iat.Unix(),
				},
			}
			// create access token in palermo service.
			sessCred, err := ctx.SessionService.Create(ctxb, sess)
			if err != nil {
				return nil, err
			}

			user := &sicily.User{
				ID:       u.Data.Id,
				Email:    u.Data.Email,
				FullName: u.Data.FullName,
				Token:    sessCred.Data.AuthToken,
			}
			return user, nil
		},
	}
}

// RefreshSession ...
func RefreshSession(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.User,
		Description: "Create user session",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"token": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"fullname": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			email, ok := params.Args["email"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			fullName, ok := params.Args["full_name"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}
			// google token
			token, ok := params.Args["token"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}

			ctxb := context.Background()
			opts := &citizens.CreateRequest{
				Data: &citizens.Citizen{
					Email:    email,
					FullName: fullName,
				},
			}
			// create user in syracusa service.
			u, err := ctx.UserService.Create(ctxb, opts)
			if err != nil {
				return nil, err
			}

			id, iat, err := idWithTime()
			if err != nil {
				return nil, err
			}
			sess := &auth.CreateRequest{
				Data: &auth.Session{
					Id:        id,
					UserId:    u.Data.Id,
					Email:     email,
					Token:     token,
					CreatedAt: iat.Unix(),
					UpdatedAt: iat.Unix(),
				},
			}
			// create access token in palermo service.
			sessCred, err := ctx.SessionService.Create(ctxb, sess)
			if err != nil {
				return nil, err
			}

			user := &sicily.User{
				ID:       u.Data.Id,
				Email:    u.Data.Email,
				FullName: u.Data.FullName,
				Token:    sessCred.Data.AuthToken,
			}

			return user, nil
		},
	}
}

// Session expose SessionMutations.
func Session(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "SessionMutations",
		Fields: graphql.Fields{
			"createSession": CreateSession(ctx),
		},
	})
}

func idWithTime() (string, *time.Time, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}

	iat := time.Now()
	id := base64.StdEncoding.EncodeToString(b)

	return id, &iat, nil
}
