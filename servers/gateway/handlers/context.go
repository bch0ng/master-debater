package handlers

import (
	"sync"

	"github.com/bch0ng/master-debater/servers/gateway/models/users"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// HandlerContext saves the required context for handlers.
type HandlerContext struct {
	Users    users.Store `json:"users"`
	CurrUser *users.User `json:"user"`
}

type WebsocketContext struct {
	Context       HandlerContext          `json:"-"`
	Connections   map[int]*websocket.Conn `json:"-"`
	Lock          *sync.Mutex             `json:"-"`
	RabbitChannel *amqp.Channel           `json:"-"`
}
