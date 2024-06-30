package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

func UnaryLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		st, _ := status.FromError(err)

		log.Info("gRPC request",
			"method", info.FullMethod,
			"duration", duration,
			"error", st.Err(),
		)

		return resp, err
	}
}

func StreamLoggingInterceptor(log *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		st, _ := status.FromError(err)

		log.Info("gRPC stream request",
			"method", info.FullMethod,
			"duration", duration,
			"error", st.Err(),
		)

		return err
	}
}
