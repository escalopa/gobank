package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/escalopa/gobank/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *GRPCServer) authenticateUser(ctx context.Context) (payload *token.Payload, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, err
	}

	// check authorization header
	authorizationHeader := md.Get(authorizationHeaderKey)
	if len(authorizationHeader) != 2 {
		return nil, fmt.Errorf("invalid authorization header format: %s", err)
	}

	// check authorization type
	if strings.ToLower(authorizationHeader[0]) != authorizationTypeBearer {
		return nil, fmt.Errorf("unsupported authentication type, provided: %s, expected: %s", authorizationHeader[0], authorizationTypeBearer)
	}

	// verify access token
	payload, err = server.tm.VerifyToken(authorizationHeader[1])
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %s", err)
	}

	return
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
