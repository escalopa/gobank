package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// MetadataKey is the key used to store the metadata in the context
	userAgent            = "user-agent"
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	xForwardForHeader    = "x-forwarded-host"
)

type Metadata struct {
	ClientIP  string
	UserAgent string
}

func (server *GRPCServer) extractMetadata(ctx context.Context) *Metadata {
	meta := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("extractMetadata: %v", md)

		if len(md[userAgent]) > 0 {
			meta.UserAgent = md[userAgent][0]
		}
		if len(md[grpcGatewayUserAgent]) > 0 {
			meta.UserAgent = md[grpcGatewayUserAgent][0]
		}
		if len(md[xForwardForHeader]) > 0 {
			meta.ClientIP = md[xForwardForHeader][0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = p.Addr.String()
	}

	return meta
}
