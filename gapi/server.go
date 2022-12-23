package gapi

import (
	"fmt"
	"log"
	"net"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/grpc/pb"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	config *util.Config
	db     db.Store
	tm     token.Maker
	pb.UnimplementedBankServiceServer
}

func NewServer(config *util.Config, store db.Store) (*GRPCServer, error) {
	maker, err := token.NewPasetoMaker(config.Get("SYMMETRIC_KEY"))
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenMaker for grpcServer, %w", err)
	}

	grpcServer := &GRPCServer{config: config, tm: maker, db: store}
	return grpcServer, nil
}

func (server *GRPCServer) Start(address string) error {
	grpcServer := grpc.NewServer()
	pb.RegisterBankServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening on %s", address)
	return nil
}
