package main

import (
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/servers"
	"github.com/escalopa/gobank/util"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

var config util.Config

func main() {
	// Load config from environment variables
	config = util.LoadConfig(".")

	// Initialize the database
	conn := db.InitDatabase(config)
	store := db.NewStore(conn)

	// Run the gRPC, Gin, and gRPC-Gateway servers concurrently
	go servers.RunGinServer(config, store)
	go servers.RunGatewayServer(config, store)
	servers.RunGRPCServer(config, store)
}
