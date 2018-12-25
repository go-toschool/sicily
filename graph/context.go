package graph

import (
	"github.com/go-toschool/palermo/auth"
	"github.com/go-toschool/platon/talks"
	"github.com/go-toschool/syracuse/citizens"
)

type Context struct {
	UserService    citizens.CitizenshipClient
	TalkService    talks.TalkingClient
	SessionService auth.AuthServiceClient
}
