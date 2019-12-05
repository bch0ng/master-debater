package main

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/db"
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/handlers"
	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/sessions"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil && CHECKRABBIT {
		log.Fatalf("%s: %s", msg, err)
	}
}

const CHECKRABBIT = false

//main is the main entry point for the server
func main() {
	// Read in ENV variables
	addr, addrExists := os.LookupEnv("ADDR")
	if !addrExists {
		addr = ":443"
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
	sessionKey, sessionKeyExists := os.LookupEnv("SESSIONKEY")
	if !sessionKeyExists {
		log.Fatalf("Environment variable SESSIONKEY not defined.")
		os.Exit(1)
	}
	redisAddr, redisAddrExists := os.LookupEnv("REDISADDR")
	if !redisAddrExists {
		log.Fatalf("Environment variable REDISADDR not defined.")
		os.Exit(1)
	}
	dsn, dsnExists := os.LookupEnv("DSN")
	if !dsnExists {
		log.Fatalf("Environment variable DSN not defined.")
		os.Exit(1)
	}
	messagesAddr, messagesAddrExists := os.LookupEnv("MESSAGESADDR")
	if !messagesAddrExists {
		log.Fatalf("Environment variable MESSAGESADDR not defined.")
		os.Exit(1)
	}
	summaryAddr, summaryAddrExists := os.LookupEnv("SUMMARYADDR")
	if !summaryAddrExists {
		log.Fatalf("Environment variable SUMMARYADDR not defined.")
		os.Exit(1)
	}
	rabbitAddr, rabbitAddrExists := os.LookupEnv("RABBITMQADDR")
	if !rabbitAddrExists {
		log.Fatalf("Environment variable RABBITMQADDR not defined.")
		os.Exit(1)
	}

	// Init sql database
	userStore, postgressErr := db.ConnectToPostgres(dsn)
	if postgressErr != nil {
		log.Fatalf("Postgress not working")
		os.Exit(1)
	}
	//Init RabbitMQ
	//"amqp://my-r:guest@localhost:5672/"
	rabbitConn, err := amqp.Dial(rabbitAddr)
	failOnError(err, "Rabbit server not available")
	ch, err := rabbitConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func() {
		fmt.Println("Rabbit MQ connection closing")
		rabbitConn.Close()
	}()

	//defer rabbitConn.Close()

	// Init Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	redisSession := sessions.NewRedisStore(redisClient, 3600)

	handlerContext := &handlers.SessionContext{
		Key:     sessionKey,
		Session: redisSession,
		User:    userStore,
	}
	websocketContext := &handlers.WebsocketContext{
		Context:       *handlerContext,
		Connections:   make(map[int]*websocket.Conn),
		Lock:          &sync.Mutex{},
		RabbitChannel: ch,
	}
	websocketContext.StartRabbitConsumer()

	// Round-robin approach for selecting
	// next address
	summaryAddrSplit := strings.Split(summaryAddr, ",")
	summaryAddrRing := ring.New(len(summaryAddrSplit))
	for _, val := range summaryAddrSplit {
		summaryAddrRing.Value = strings.TrimSpace(val)
		summaryAddrRing = summaryAddrRing.Next()
	}
	messagesAddrSplit := strings.Split(messagesAddr, ",")
	messagesAddrRing := ring.New(len(messagesAddrSplit))
	for _, val := range messagesAddrSplit {
		messagesAddrRing.Value = strings.TrimSpace(val)
		messagesAddrRing = messagesAddrRing.Next()
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./summary/"))
	mux.Handle("/summary/", http.StripPrefix("/summary", fileServer))

	// Reverse proxies
	mux.HandleFunc("/v1/summary", func(w http.ResponseWriter, r *http.Request) {

		var sessionState *handlers.SessionState = new(handlers.SessionState)
		_, err := sessions.GetState(r, handlerContext.Key, redisSession, sessionState)
		if err != nil {
			log.Fatal(err)
		}
		if sessionState.User != nil {
			jsonObj, err := json.Marshal(sessionState.User)

			if err != nil {
				log.Fatal(err)
			}
			r.Header.Set("X-User", string(jsonObj))
		} else {
			r.Header.Del("X-User")
		}
		summaryProxyAddr := summaryAddrRing.Value.(string)
		log.Println(summaryProxyAddr)
		summaryProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: summaryProxyAddr})
		summaryProxy.ServeHTTP(w, r)
		summaryAddrRing = summaryAddrRing.Next()
	})
	mux.HandleFunc("/v1/channels", func(w http.ResponseWriter, r *http.Request) {
		var sessionState *handlers.SessionState = new(handlers.SessionState)
		_, err := sessions.GetState(r, handlerContext.Key, redisSession, sessionState)
		if err != nil {
			log.Fatal(err)
		}
		if sessionState.User != nil {
			json, err := json.Marshal(sessionState.User)
			if err != nil {
				log.Fatal(err)
			}
			r.Header.Set("X-User", string(json))
		} else {
			r.Header.Del("X-User")
		}
		channelsProxyAddr := messagesAddrRing.Value.(string)
		channelsProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: channelsProxyAddr})
		channelsProxy.ServeHTTP(w, r)
		messagesAddrRing = messagesAddrRing.Next()
	})
	mux.HandleFunc("/v1/channels/", func(w http.ResponseWriter, r *http.Request) {
		var sessionState *handlers.SessionState = new(handlers.SessionState)
		_, err := sessions.GetState(r, handlerContext.Key, redisSession, sessionState)
		if err != nil {
			log.Fatal(err)
		}
		if sessionState.User != nil {
			json, err := json.Marshal(sessionState.User)
			if err != nil {
				log.Fatal(err)
			}
			r.Header.Set("X-User", string(json))
		} else {
			r.Header.Del("X-User")
		}
		channelsProxyAddr := messagesAddrRing.Value.(string)
		channelsProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: channelsProxyAddr})
		channelsProxy.ServeHTTP(w, r)
		messagesAddrRing = messagesAddrRing.Next()
	})
	mux.HandleFunc("/v1/messages/", func(w http.ResponseWriter, r *http.Request) {
		var sessionState *handlers.SessionState = new(handlers.SessionState)
		_, err := sessions.GetState(r, handlerContext.Key, redisSession, sessionState)
		if err != nil {
			log.Fatal(err)
		}
		if sessionState.User != nil {
			json, err := json.Marshal(sessionState.User)
			if err != nil {
				log.Fatal(err)
			}
			r.Header.Set("X-User", string(json))
		} else {
			r.Header.Del("X-User")
		}
		messagesProxyAddr := messagesAddrRing.Value.(string)
		messagesProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: messagesProxyAddr})
		messagesProxy.ServeHTTP(w, r)
		messagesAddrRing = messagesAddrRing.Next()
	})

	mux.HandleFunc("/v2/openchannels", func(w http.ResponseWriter, r *http.Request) {
		/*
			var sessionState *handlers.SessionState = new(handlers.SessionState)
			_, err := sessions.GetState(r, handlerContext.Key, redisSession, sessionState)
			if err != nil {
				log.Fatal(err)
			}
			if sessionState.User != nil {
				json, err := json.Marshal(sessionState.User)
				if err != nil {
					log.Fatal(err)
				}
				r.Header.Set("X-User", string(json))
			} else {
				r.Header.Del("X-User")
			}
		*/
		fmt.Println("v2 get all open channels called")
		channelsProxyAddr := messagesAddrRing.Value.(string)
		channelsProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: channelsProxyAddr})
		channelsProxy.ServeHTTP(w, r)
		messagesAddrRing = messagesAddrRing.Next()
	})

	fmt.Println("Sanity Check:Version wedsMorning1awef")

	mux.HandleFunc("/v1/users", handlerContext.UsersHandler)
	mux.HandleFunc("/v1/users/", handlerContext.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", handlerContext.SessionsHandler) //Sessions handler handles login attempts
	mux.HandleFunc("/v1/sessions/", handlerContext.SpecificSessionHandler)
	mux.HandleFunc("/ws", websocketContext.WebSocketHandler)
	corsMux := handlers.NewCORSMiddleware(mux)
	log.Printf("Server is open and listening on %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsMux))
}
