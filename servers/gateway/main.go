package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"

	"github.com/bch0ng/master-debater/servers/gateway/db"
	"github.com/bch0ng/master-debater/servers/gateway/handlers"
	"github.com/bch0ng/master-debater/servers/gateway/models/users"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

//main is the main entry point for the server
func main() {
	// Read in ENV variables
	addr, addrExists := os.LookupEnv("ADDR")
	if !addrExists {
		addr = ":443"
	}

	// JWT Key
	_, jwtSecretExists := os.LookupEnv("JWT_SECRET")
	if !jwtSecretExists {
		log.Fatalf("Environment variable JWT_SECRET not defined.")
		os.Exit(1)
	}

	tlsCertPath, tlsCertExists := os.LookupEnv("TLSCERT")
	if !tlsCertExists {
		log.Fatalf("Environment variable TLSCERT not defined.")
		os.Exit(1)
	}
	tlsKeyPath, tlsKeyExists := os.LookupEnv("TLSKEY")
	if !tlsKeyExists {
		log.Fatalf("Environment variable TLSKEY not defined.")
		os.Exit(1)
	}

	// MySQL conenction
	dsn, dsnExists := os.LookupEnv("DSN")
	if !dsnExists {
		log.Fatalf("Environment variable DSN not defined.")
		os.Exit(1)
	}
	MySQLStore, err := db.ConnectToPostgres(dsn)
	if err != nil {
		log.Fatalf("MySQL not working")
		os.Exit(1)
	}

	// Microservice address
	microAddr, microAddrExists := os.LookupEnv("MICRO_ADDR")
	if !microAddrExists {
		log.Fatalf("Environment variable MICRO_ADDR not defined.")
		os.Exit(1)
	}

	handlerContext := &handlers.HandlerContext{
		Users:    MySQLStore,
		CurrUser: new(users.User),
	}

	// RabbitMQ address
	rabbitmqAddr, rabbitmqAddrExists := os.LookupEnv("RABBITMQ_ADDR")
	if !rabbitmqAddrExists {
		log.Fatalf("Environment variable RABBITMQ_ADDR not defined.")
		os.Exit(1)
	}
	rabbitConn, err := amqp.Dial(rabbitmqAddr)
	if err != nil {
		log.Fatal("Rabbit server not available")
	}
	ch, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
	}
	defer func() {
		fmt.Println("Rabbit MQ connection closing")
		rabbitConn.Close()
	}()
	websocketContext := &handlers.WebsocketContext{
		Context:       *handlerContext,
		Connections:   make(map[int]*websocket.Conn),
		Lock:          &sync.Mutex{},
		RabbitChannel: ch,
	}
	websocketContext.StartRabbitConsumer()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/debate/", func(w http.ResponseWriter, r *http.Request) {
		if handlerContext.CurrUser.Username != "" {
			json, err := json.Marshal(handlerContext.CurrUser)
			if err != nil {
				log.Fatal(err)
			}
			r.Header.Set("X-User", string(json))
		} else {
			r.Header.Del("X-User")
		}
		channelsProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: microAddr})
		channelsProxy.ServeHTTP(w, r)
	})

	mux.HandleFunc("/api/user/create", handlerContext.CreateUserHandler)
	mux.HandleFunc("/api/user/login", handlerContext.LoginUserHandler)

	mux.Handle("/api/chatroom", handlers.NewJWTMiddleware(testRoute(mux)))
	mux.Handle("/api/user/logout", handlers.NewJWTMiddleware(testRoute2(mux, handlerContext)))

	corsMux := handlers.NewCORSMiddleware(mux)
	log.Printf("Server is open and listening on %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsMux))
}

func testRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO WORLD"))
	})
}

func testRoute2(h http.Handler, context *handlers.HandlerContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.LogoutUserHandler(w, r)
	})
}
