package main

import (
	"context"
	"log"
	"net"
	"net/http"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/gapi"
	"github.com/escalopa/gobank/grpc/pb"
	"github.com/escalopa/gobank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
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

	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatalf("cannot listen on port 8000, err: %s", err)
	}

	log.Printf("Gateway server is listening on 8000")
	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("cannot start HTTP server, err: %s", err)
	}
}
