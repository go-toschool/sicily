package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-toschool/palermo"
	"github.com/go-toschool/palermo/auth"
	"github.com/go-toschool/sicily"

	"google.golang.org/grpc"

	"github.com/go-toschool/sicily/cmd/server/session"
	"github.com/go-toschool/sicily/cmd/server/users"
	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/mutation"
	"github.com/go-toschool/sicily/graph/queries"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/graphql-go/graphql"
)

type stringValueKey string

const (
	// AuthUserIDContextKey key for context
	AuthUserIDContextKey sicily.StringValueKey = "userIdContextKey"
	// ContentTypeGraphQL ...
	ContentTypeGraphQL = "application/graphql"

	tokenTypePrefix     = "Bearer "
	tokenHeaderKey      = "Authorization"
	authTokenCookieName = "access_token"
)

func main() {
	// Connect services
	citizensConn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	check("citizens connection:", err)

	// platonConn, err := grpc.Dial("localhost:8002", grpc.WithInsecure())
	// check("platon connection:", err)

	palermoConn, err := grpc.Dial("localhost:8003", grpc.WithInsecure())
	check("palermo connection:", err)
	// Initialize citizen client
	citizenSvc := citizens.NewCitizenshipClient(citizensConn)
	// talksSvc := talks.NewTalkingClient(platonConn)
	palermoSvc := auth.NewAuthServiceClient(palermoConn)

	graphCtx := &graph.Context{
		UserService: citizenSvc,
		// TalkService:    talksSvc,
		SessionService: palermoSvc,
	}

	// user schema
	userSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queries.Users(graphCtx),
		Mutation: mutation.Users(graphCtx),
	})
	check("user schema:", err)

	uc := &users.Context{
		User:   citizenSvc,
		Schema: userSchema,
	}
	uac := &AuthCors{
		User:           citizenSvc,
		SessionService: palermoSvc,
		Next:           uc.Handle(users.Users),
	}
	// private endpoint
	http.HandleFunc("/users", uac.CheckCorsAndAuth())

	// // talk schema
	// talkSchema, err := graphql.NewSchema(graphql.SchemaConfig{
	// 	Query:    queries.Talks(graphCtx),
	// 	Mutation: mutation.Talks(graphCtx),
	// })
	// check("talk schema:", err)

	// tc := &apiTalks.Context{
	// 	Talk:   talksSvc,
	// 	Schema: talkSchema,
	// }
	// tac := ac.AddHandler(tc.Handle(apiTalks.Talks))
	// // private endpoint
	// http.HandleFunc("/tasks", tac.CheckCorsAndAuth())

	// session
	sessionSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queries.Session(graphCtx),
		Mutation: mutation.Session(graphCtx),
	})
	check("session schema:", err)

	// public endpoints
	sc := &session.Context{
		User:   citizenSvc,
		Schema: sessionSchema,
	}
	sac := &AuthCors{
		User:           citizenSvc,
		SessionService: palermoSvc,
		Next:           sc.Handle(session.Session),
	}
	http.HandleFunc("/session", sac.CheckCors())

	log.Println("Now server is running on port 3000")
	http.ListenAndServe(":3000", nil)
}

// AuthCors ...
type AuthCors struct {
	User           citizens.CitizenshipClient
	SessionService auth.AuthServiceClient
	Next           http.Handler
}

// CheckCors ...
func (ac *AuthCors) CheckCors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		ac.Next.ServeHTTP(w, r)
	}
}

// CheckCorsAndAuth ...
func (ac *AuthCors) CheckCorsAndAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		cred, err := parseAuthCredentials(r)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		session, err := ac.SessionService.Get(ctx, &auth.GetRequest{
			Data: &auth.SessionCredentials{
				ValidationToken: cred.ValidationToken,
				AuthToken:       cred.AuthToken,
			},
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx1 := r.Context()
		ctx1 = setUserIDToRequestContext(ctx1, session.Data.UserId)
		ac.Next.ServeHTTP(w, r.WithContext(ctx1))
	}
}

func setUserIDToRequestContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, AuthUserIDContextKey, userID)
}

func parseAuthToken(r *http.Request) (string, error) {
	header := r.Header.Get(tokenHeaderKey)
	if !strings.HasPrefix(header, tokenTypePrefix) {
		return "", errors.New("auth: no token authorization header present")
	}
	return header[len(tokenTypePrefix):], nil
}

func parseAuthCredentials(r *http.Request) (*palermo.SessionCredentials, error) {
	authToken, err := parseAuthToken(r)
	if err != nil {
		return nil, err
	}
	validationToken, err := parseValidationToken(r)
	if err != nil {
		return nil, err
	}

	return &palermo.SessionCredentials{
		AuthToken:       authToken,
		ValidationToken: validationToken,
	}, nil
}

func parseValidationToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(authTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func check(section string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s %v", section, err))
	}
}
