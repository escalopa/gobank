package servers

import (
	"context"
	"log"
	"net"
	"net/http"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/gapi"
	"github.com/escalopa/gobank/pb"
	"github.com/escalopa/gobank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func RunGatewayServer(config util.Config, store db.Store) {
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
