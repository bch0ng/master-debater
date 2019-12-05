package handlers

import (
	"sync"

	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/sessions"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// SessionContext saves the required context for handlers.
type SessionContext struct {
	Key     string               `json:"-"`
	Session *sessions.RedisStore `json:"session"`
	User    users.Store          `json:"user"`
}
type WebsocketContext struct {
	Context       SessionContext          `json:"-"`
	Connections   map[int]*websocket.Conn `json:"-"`
	Lock          *sync.Mutex             `json:"-"`
	RabbitChannel *amqp.Channel           `json:"-"`
}
