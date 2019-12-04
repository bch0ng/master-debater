package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bch0ng/master-debater/servers/gateway/db"
	"github.com/bch0ng/master-debater/servers/gateway/handlers"
)

func failOnError(err error, msg string) {
	if err != nil && CHECKRABBIT {
		log.Fatalf("%s: %s", msg, err)
	}
}

const CHECKRABBIT = true

//main is the main entry point for the server
func main() {
	// Read in ENV variables
	addr, addrExists := os.LookupEnv("ADDR")
	if !addrExists {
		addr = ":3003"
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

	handlerContext := &handlers.HandlerContext{
		User: MySQLStore,
	}

	mux.HandleFunc("/api/user/create", handlerContext.CreateUserHandler)
	mux.HandleFunc("/api/user/login", handlerContext.LoginUserHandler)

	mux.HandleFunc("/api/chatroom", handlers.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write("HELLO WORLD")
	}))

	corsMux := handlers.NewCORSMiddleware(mux)
	log.Printf("Server is open and listening on %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsMux))
}
