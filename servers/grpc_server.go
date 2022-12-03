package servers

import (
	"log"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/gapi"
	"github.com/escalopa/gobank/util"
)

func RunGRPCServer(config util.Config, store db.Store) {
	grpcServer, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPC server, err: %s", err)
	}

	log.Printf("GRPC server is listening on %s", config.GRPCPort)
	if err := grpcServer.Start(config.GRPCPort); err != nil {
		log.Fatalf("cannot start gRPC server address: %s, err: %s", config.GRPCPort, err)
	}
}
