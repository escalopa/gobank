package gapi

import (
	"context"
	"database/sql"
	"strings"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/pb"
	"github.com/escalopa/gobank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GRPCServer) UpdateUser(ctx context.Context, req *pb.UserUpdateRequest) (*pb.UserResponse, error) {
	// Authenticate user
	payload, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	// Get user from database
	user, err := server.getUser(ctx, payload.Username)
	if err != nil {
		return nil, err
	}

	var arg db.UpdateUserParams

	// Add email if provided
	// TODO: validate email
	if req.GetEmail() != "" {
		if req.GetEmail() == user.Email {
			return nil, status.Error(codes.AlreadyExists, "new email cannot be equal to current email")
		}
		arg.Email = sql.NullString{
			String: req.GetEmail(),
			Valid:  true,
		}
	}

	// Add fullname if provided
	// TODO: validate fullname (no numbers, only letters & spaces)
	if req.GetFullName() != "" {

		// Remove all redundant spaces between the names
		formattedFullName := strings.Join(strings.Fields(req.GetFullName()), " ")
		if formattedFullName == user.FullName {
			return nil, status.Errorf(codes.AlreadyExists, "new full_name cannot be equal to current full_name")
		}
		arg.FullName = sql.NullString{
			String: formattedFullName,
			Valid:  true,
		}
	}

	// Add new password if provided
	if req.GetPassword() != nil {
		if err = util.CheckHashedPassword(user.HashedPassword, req.GetPassword().GetOldPassword()); err == nil {
			return nil, status.Errorf(codes.AlreadyExists, "new password cannot be equal to current password")
		}

		hashedPassword, err := util.GenerateHashPassword(req.GetPassword().GetNewPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.Email = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	res := fromDBUserToPbUserResponse(user)
	return res, nil
}
