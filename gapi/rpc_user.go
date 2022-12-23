package gapi

import (
	"context"
	"database/sql"
	"strings"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/grpc/pb"
	"github.com/escalopa/gobank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Get User from DB by Username
	user, err := server.getUser(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	// Check User's Password
	err = util.CheckHashedPassword(user.HashedPassword, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := server.tm.CreateToken(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	// Generate New Access Token for User
	refreshToken, refreshPayload, err := server.tm.CreateRefreshToken(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}

	// Create new session for User
	md := server.extractMetadata(ctx)
	session, err := server.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.GetUsername(),
		RefreshToken: refreshToken,
		UserAgent:    md.UserAgent,
		ClientIp:     md.ClientIP,
		ExpiresAt:    refreshPayload.ExpireAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	res := &pb.LoginResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpireAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpireAt),
		User:                  fromDBUserToPbUserResponse(user),
	}
	return res, nil
}

func (server *GRPCServer) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	hashPassword, err := util.GenerateHashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot hash password")
	}

	user, err := server.db.CreateUser(ctx, db.CreateUserParams{
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
	user, err := server.db.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, status.Errorf(codes.NotFound, "user %s not found", username)
		}
		return db.User{}, status.Errorf(codes.Internal, "cannot get user %s: %v", username, err)
	}
	return user, nil
}

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

	user, err = server.db.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	res := fromDBUserToPbUserResponse(user)
	return res, nil
}
