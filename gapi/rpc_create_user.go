package gapi

import (
	"context"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/pb"
	"github.com/escalopa/gobank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GRPCServer) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	hashPassword, err := util.GenerateHashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot hash password")
	}

	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username %s already exists", req.GetUsername())
			}
		}
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	res := fromDBUserToPbUserResponse(user)
	return res, nil
}
