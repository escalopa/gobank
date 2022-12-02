package main

import (
	"log"

	"github.com/escalopa/gobank/api"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/gapi"
	"github.com/escalopa/gobank/util"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
)

var config util.Config

func main() {
	var err error

	config, err = util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot read configuration", err)
	}

	conn, err := util.InitDatabase(config)
	if err != nil {
		log.Fatalf("cannot open connection to db, err: %s", err)
	}

	store := db.NewStore(conn)

	// runGinServer(config, store)
	runGRPCServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	ginServer, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create HTTP server, err: %s", err)
	}

	if err := ginServer.Start(config.HTTPPort); err != nil {
		log.Fatalf("cannot start server address: %s, err: %s", config.HTTPPort, err)
	}
}

func runGRPCServer(config util.Config, store db.Store) {
	grpcServer, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPC server, err: %s", err)
	}

	// Init gRPC server
	if err := grpcServer.Start(config.GRPCPort); err != nil {
		log.Fatalf("cannot start gRPC server address: %s, err: %s", config.GRPCPort, err)
	}
}
