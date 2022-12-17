package gapi

import (
	"context"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/grpc/pb"
	"github.com/escalopa/gobank/util"
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
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	// Generate New Access Token for User
	refreshToken, refreshPayload, err := server.tokenMaker.CreateRefreshToken(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}

	// Create new session for User
	md := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
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
