package main

import (
	"log"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/gapi"
	"github.com/escalopa/gobank/util"
)

func main() {
	// Load config from environment variables
	config := util.NewConfig()

	// Initialize the database
	conn := db.InitDatabase(config)
	store := db.NewStore(conn)
	grpcServer, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPC server, err: %s", err)
	}

	log.Printf("GRPC server is listening on %s", "8000")
	if err := grpcServer.Start("0.0.0.0:8000"); err != nil {
		log.Fatalf("cannot start gRPC server address: %s, err: %s", "8000", err)
	}
}
