package gapi

import (
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func fromDBUserToPbUserResponse(user db.User) *pb.UserResponse {
	return &pb.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
