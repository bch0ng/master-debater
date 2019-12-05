package main

import (
	"log"
	"net/http"
	"os"
)

//main is the main entry point for the server
func main() {
	addr, addrExists := os.LookupEnv("ADDR")
	if !addrExists {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("../../clients/summary"))
	mux.Handle("/summary/", http.StripPrefix("/summary", fileServer))
	mux.HandleFunc("/v1/summary", SummaryHandler)
	log.Printf("Server is open and listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
