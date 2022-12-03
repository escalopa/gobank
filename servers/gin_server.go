package servers

import (
	"log"

	"github.com/escalopa/gobank/api"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/util"
)

func RunGinServer(config util.Config, store db.Store) {
	ginServer, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create HTTP server, err: %s", err)
	}

	log.Printf("HTTP server is listening on %s", config.HTTPPort)
	if err := ginServer.Start(config.HTTPPort); err != nil {
		log.Fatalf("cannot start server address: %s, err: %s", config.HTTPPort, err)
	}
}
