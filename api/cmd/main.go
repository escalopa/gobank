package main

import (
	"log"

	"github.com/escalopa/gobank/api/handlers"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/util"
)

func main() {
	// Load config from environment variables
	config := util.NewConfig()

	// Initialize the database
	conn := db.InitDatabase(config)
	store := db.NewStore(conn)

	ginServer, err := handlers.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create HTTP server, err: %s", err)
	}

	log.Printf("HTTP server is listening on %s", "8000")
	if err := ginServer.Start("0.0.0.0:8000"); err != nil {
		log.Fatalf("cannot start server address: %s, err: %s", "8000", err)
	}
}
