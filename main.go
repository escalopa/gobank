package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/escalopa/gobank/api"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/gapi"
	"github.com/escalopa/gobank/pb"
	"github.com/escalopa/gobank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
	"google.golang.org/protobuf/encoding/protojson"
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

	go runGinServer(config, store)
	go runGatewayServer(config, store)
	runGRPCServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	ginServer, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create HTTP server, err: %s", err)
	}

	log.Printf("HTTP server is listening on %s", config.HTTPPort)
	if err := ginServer.Start(config.HTTPPort); err != nil {
		log.Fatalf("cannot start server address: %s, err: %s", config.HTTPPort, err)
	}
}

func runGRPCServer(config util.Config, store db.Store) {
	grpcServer, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPC server, err: %s", err)
	}

	log.Printf("GRPC server is listening on %s", config.GRPCPort)
	if err := grpcServer.Start(config.GRPCPort); err != nil {
		log.Fatalf("cannot start gRPC server address: %s, err: %s", config.GRPCPort, err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	grpcServer, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPC server, err: %s", err)
	}

	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOpts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := pb.RegisterBankServiceHandlerServer(ctx, grpcMux, grpcServer); err != nil {
		log.Fatalf("cannot register gRPC server, err %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.GatewayPort)
	if err != nil {
		log.Fatalf("cannot listen on port %s, err: %s", config.GatewayPort, err)
	}

	log.Printf("Gateway server is listening on %s", config.GatewayPort)
	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("cannot start HTTP server, err: %s", err)
	}

}
