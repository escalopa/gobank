package gapi

import (
	"context"
	"database/sql"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GRPCServer) GetUser(ctx context.Context, req *pb.Username) (*pb.UserResponse, error) {
	payload, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if req.GetUsername() != payload.Username {
		return nil, status.Error(codes.Unauthenticated, "requested username doesn't match the provided in token")
	}

	user, err := server.getUser(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	res := fromDBUserToPbUserResponse(user)
	return res, nil
}

func (server *GRPCServer) getUser(ctx context.Context, username string) (db.User, error) {
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, status.Errorf(codes.NotFound, "user %s not found", username)
		}
		return db.User{}, status.Errorf(codes.Internal, "cannot get user %s: %v", username, err)
	}
	return user, nil
}
