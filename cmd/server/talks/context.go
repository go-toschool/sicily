package talks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-toschool/platon/talks"
	"github.com/go-toschool/sicily"
	"github.com/graphql-go/graphql"
)

// GraphRequest struct to unmarshal query.
type GraphRequest struct {
	Query string `json:"query"`
}

// Context ...
type Context struct {
	Talk   talks.TalkingClient
	Schema graphql.Schema
}

// Handle creates a new bounded Handler with context.
func (c *Context) Handle(h HandlerFunc) *Handler {
	return &Handler{c, h}
}

// ExecuteQuery ...
func (c *Context) ExecuteQuery(query, userID string) *graphql.Result {
	ctx := context.WithValue(context.Background(), sicily.UserIDKey, userID)
	result := graphql.Do(graphql.Params{
		Schema:        c.Schema,
		RequestString: query,
		Context:       ctx,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v\n", result.Errors)
	}

	return result
}

// HandlerFunc function handler signature used by sigiriya application.
type HandlerFunc func(*Context, http.ResponseWriter, *http.Request)

// Handler ...
type Handler struct {
	ctx    *Context
	handle HandlerFunc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handle(h.ctx, w, r)
}
