package handlers

import (
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
)

// HandlerContext saves the required context for handlers.
type HandlerContext struct {
	Users users.Store `json:"users"`
}
