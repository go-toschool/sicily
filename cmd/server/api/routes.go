package api

import (
	"github.com/go-toschool/sicily"
	"github.com/go-toschool/sicily/cmd/server/firewall"
	"github.com/gorilla/mux"
)

const (
	// AuthUserIDContextKey key for context
	AuthUserIDContextKey sicily.StringValueKey = "userIdContextKey"
	authTokenCookieName                        = "access_token"
)

func Routes(ctx *Context) *mux.Router {
	r := mux.NewRouter()

	firewall := firewall.NewAuth(ctx.Session)
	api := ctx.Handle(API)
	r.HandleFunc("/graphql", firewall.CheckCorsAndToken(api))

	return r
}
