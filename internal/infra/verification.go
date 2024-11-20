package infra

import (
	"context"

	"github.com/s21platform/society-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Verifcation(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	userIDs := md["uuid"]
	if !ok || len(userIDs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no uuid found in metadata")
	}
	ctx = context.WithValue(ctx, config.KeyUUID, userIDs[0])
	return handler(ctx, req)
}
