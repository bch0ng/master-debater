package handlers

import (
	"time"

	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
)

// SessionState logs a user's session.
type SessionState struct {
	StartTime time.Time   `json:"time"`
	User      *users.User `json:"user"`
}
