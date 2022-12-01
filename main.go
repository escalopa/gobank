package main

import (
	"log"

	"github.com/escalopa/go-bank/api"
	db "github.com/escalopa/go-bank/db/sqlc"
	"github.com/escalopa/go-bank/util"
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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server, err: %s", err)
	}

	if err := server.Start(config.Port); err != nil {
		log.Fatalf("cannot start server address: %s, err: %s", config.Port, err)
	}

}
