package infra

import (
	"context"
	"github.com/s21platform/society-service/internal/config"

	"google.golang.org/grpc"

	logger_lib "github.com/s21platform/logger-lib"
)

func Logger(logger *logger_lib.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = context.WithValue(ctx, config.KeyLogger, logger)
		return handler(ctx, req)
	}
}
