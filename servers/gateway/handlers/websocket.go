package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	//"time"
	//"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
)

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutinews.

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket
const CHECK_ORIGIN = false

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// This function's purpose is to reject websocket upgrade requests if the
		// origin of the websockete handshake request is coming from unknown domainwsc.
		// This prevents some random domain from opening up a socket with your server.
		// TODO: make sure you modify this for your HW to check if r.Origin is your host

		origin := r.Header.Get("Origin")
		fmt.Println("Upgrade attempt, origin:", origin)
		return !CHECK_ORIGIN || (origin != "https://www.sumsumsummary.me")
	},
}

// Control messages for websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// This is a struct to read our message into
type msg struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
}
type DebateMessage struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
	ChannelID   int    `json:"channelId"`
	Handle      string `json:"handle"`
	Username    string `json:"username"`
}
type Creator struct {
	id int `json:"id"`
}
type Channel struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	Members     []int   `json:"members"`
	CreatedAt   string  `json:"createdAt"`
	Creator     Creator `json:"creator"`
	EditedAt    string  `json:"editedAt"`
}
type RabbitMessage struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	MessageID string `json:"messageID"`
	ChannelID string `json:"channelID"`
	Channel   string `json:"channel"`
	UserIDs   []int  `json:"userIDs"`
}

const DEBATE_MESSAGE = "debateMessage"

func (wsc *WebsocketContext) CheckUserAuth(r *http.Request, authToken string) error {
	/*Validate the authtoken, if the token validates then the user logged in previously.
	 */
	if wsc.Context.CurrUser.Username == "" {
		return errors.New("user is not authenticated")
	}
	fmt.Println("Socket auth validated:", authToken)
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (wsc *WebsocketContext) StartRabbitConsumer() {
	//We declare the queue anyways just to make sure it exists.
	//If it already exists nothing will happen
	rabbitQueue, err := wsc.RabbitChannel.QueueDeclare(
		"debate", // name
		true,     // durable
		false,    // delete when usused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")
	//Consume returns a (<-chan Delivery, error)
	msgs, err := wsc.RabbitChannel.Consume(
		rabbitQueue.Name, // queue
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			m := RabbitMessage{}
			err := json.Unmarshal([]byte(d.Body), &m)
			if err != nil {
				fmt.Println("Error reading json.", err)
				break
			}
			if m.UserIDs != nil {
				wsc.ReplyToConnections(m.UserIDs, m)
			} else {
				wsc.ReplyToAllConnectionsJson(m)
			}
			//wsc.HandleRabbitMessage(m)
		}
	}()
}

//Redundant for now.
func (wsc *WebsocketContext) HandleRabbitMessage(m RabbitMessage) {
	messageType := m.Type
	if messageType == "channel-new" {
	} else if messageType == "channel-update" {
	} else if messageType == "channel-delete" {
	} else if messageType == "message-new" {

	} else if messageType == "message-update" {

	} else if messageType == "message-delete" {

	} else {
		fmt.Println("Invalid message type:", messageType)
	}
}

// UsersHandler handles all requests for general user actions
func (wsc *WebsocketContext) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	fmt.Println("Connection attempt, Origin:", origin)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Connection error 403")
		http.Error(w, "Websocket Connection Refused", 403)
		return
	} else {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			fmt.Println("No auth in header, Attempting to get it from queryString")
			auth = r.FormValue("auth")
		}
		if auth == "" {
			fmt.Println("No auth in header or queryString, error 403")
			http.Error(w, "Websocket Connection Refused No Auth", 403)
		} else {
			err := wsc.CheckUserAuth(r, auth)
			if err != nil {
				http.Error(w, "Websocket Connection Refused, Auth Bad", 403)
			} else {
				connId := wsc.InsertConnection(conn)
				go wsc.echoSocketLoop(conn, connId)
			}
		}
	}
}

// Thread-safe method for inserting a connection
func (wsc *WebsocketContext) InsertConnection(conn *websocket.Conn) int {
	wsc.Lock.Lock()
	connId := len(wsc.Connections)
	// insert socket connection
	//wsc.Connections = append(wsc.Connections, conn)
	wsc.Connections[connId] = conn
	wsc.Lock.Unlock()
	return connId
}

// Thread-safe method for inserting a connection
func (wsc *WebsocketContext) RemoveConnection(connId int) {
	wsc.Lock.Lock()
	// insert socket connection
	//wsc.Connections = append(wsc.Connections[:connId], wsc.Connections[connId+1:]...)
	delete(wsc.Connections, connId)
	wsc.Lock.Unlock()
}

//Takes a object and sends back a json string
func (wsc *WebsocketContext) WriteJSONReply(conn *websocket.Conn, message DebateMessage) error {
	return conn.WriteJSON(message)
}
func (wsc *WebsocketContext) ForwardDebateMessage(conn *websocket.Conn, message DebateMessage) error {
	return wsc.ReplyToAllConnectionsJson(message)
}
func (wsc *WebsocketContext) SendJSONReply(id int, message DebateMessage) error {
	return wsc.Connections[id].WriteJSON(message)
}

func (wsc *WebsocketContext) ReplyToConnections(ids []int, messageObj interface{}) error {
	var writeError error
	for _, num := range ids {
		writeError = wsc.Connections[num].WriteJSON(messageObj)
		if writeError != nil {
			return writeError
		}
	}
	return nil
}

//Takes a obj and replies to all connections as json
func (wsc *WebsocketContext) ReplyToAllConnectionsJson(messageObj interface{}) error {
	var writeError error
	for _, conn := range wsc.Connections {
		writeError = conn.WriteJSON(messageObj)
		if writeError != nil {
			return writeError
		}
	}
	return nil
}

//Takes a message as text and replies to all connections
func (wsc *WebsocketContext) ReplyToAllConnections(messageType int, data []byte) error {
	var writeError error
	for _, conn := range wsc.Connections {
		writeError = conn.WriteMessage(messageType, data)
		if writeError != nil {
			return writeError
		}
	}
	return nil
}

func (wsc *WebsocketContext) EmptySocketLoop(conn *websocket.Conn, connId int) {
	defer conn.Close()
	defer wsc.RemoveConnection(connId)
	for {
		fmt.Println("New Empty SocketLoop")
		msgType, connMessage, err := conn.ReadMessage()
		if msgType == CloseMessage {
			fmt.Println("Close message received.")
			break
		} else if msgType == TextMessage || msgType == BinaryMessage {
			m := DebateMessage{}
			err := json.Unmarshal([]byte(connMessage), &m)
			if err != nil {
				fmt.Println("Error reading json.", err)
				conn.Close()
				break
			}
			fmt.Printf("Got message: %#v\n", m)
			if m.MessageType == DEBATE_MESSAGE {
				err = wsc.ForwardDebateMessage(conn, m)
				if err != nil {
					fmt.Println("Error writing debate response to socket:", err)
					break
				}
			} else {
				err = wsc.WriteJSONReply(conn, m)
				if err != nil {
					fmt.Println("Error writing response to socket:", err)
					break
				}
			}
		} else if err != nil {
			fmt.Println("Error read websocketMessage:", err)
			break
		}
	}
	fmt.Println("Connection Loop Broken:")
}
func (wsc *WebsocketContext) echoSocketLoop(conn *websocket.Conn, connId int) {
	defer conn.Close()
	defer wsc.RemoveConnection(connId)
	for {
		fmt.Println("New Echo SocketLoop")
		msgType, connMessage, err := conn.ReadMessage()
		if msgType == CloseMessage {
			fmt.Println("Close message received.")
			break
		} else if msgType == TextMessage || msgType == BinaryMessage {
			m := DebateMessage{}
			//err := conn.ReadJSON(&m)  //This causes halting.
			err := json.Unmarshal([]byte(connMessage), &m)
			if err != nil {
				fmt.Println("Error reading json.", err)
				conn.Close()
				break
			}
			fmt.Printf("Got message: %#v\n", m)
			err = wsc.WriteJSONReply(conn, m)
			if err != nil {
				fmt.Println("Error writing response to socket:", err)
				break
			}
			err = wsc.ReplyToAllConnections(TextMessage, append([]byte("Server Reply Append: "), connMessage...))
			if err != nil {
				fmt.Println("Message writing response to all sockets:", err)
				break
			}
		} else if err != nil {
			fmt.Println("Error read websocketMessage:", err)
			break
		}
	}
	fmt.Println("Connection Loop Broken:")
}
