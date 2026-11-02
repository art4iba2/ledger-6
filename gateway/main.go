package main

import (
	"gateway/internal/api"
	"ledger"
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           api.NewHandler(ledger.NewLedger()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Gateway listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
