package main

import (
	"log"

	"github.com/escalopa/gobank/api/handlers"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/util"
)

//	@title			Gobank API
//	@version		1.0
//	@description	Gobank is a SAAP that allows users to create accounts and transfer money between them.
//
//	@contact.email	ahmad.helaly.dev@gmail.com
//	@contact.name	Ahmad Helaly

//	@securityDefinitions.apikey	bearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer <token>
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
