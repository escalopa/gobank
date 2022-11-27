package main

import (
	"database/sql"
	"log"

	"github.com/escalopa/go-bank/api"
	db "github.com/escalopa/go-bank/db/sqlc"
	"github.com/escalopa/go-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot read configuration", err)
	}

	conn, err := sql.Open(config.DB.Driver, config.DB.ConnectionString)
	if err != nil {
		log.Fatalf("cannot open connection to db err: %s", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.App.Port); err != nil {
		log.Fatalf("cannot start server address: %s, err: %s", config.App.Port, err)
	}

}
