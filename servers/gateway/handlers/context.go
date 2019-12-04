package handlers

import (
	"github.com/bch0ng/master-debater/servers/gateway/models/users"
	"github.com/bch0ng/master-debater/servers/gateway/sessions"
)

// HandlerContext saves the required context for handlers.
type HandlerContext struct {
	Users     users.Store          `json:"users"`
	Blacklist *sessions.RedisStore `json:"blacklist"`
}
