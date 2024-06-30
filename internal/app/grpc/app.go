package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "vizapSSO/internal/grpc/auth"
	"vizapSSO/internal/interceptor"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, GRPCPort int) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryLoggingInterceptor(log)),
		grpc.StreamInterceptor(interceptor.StreamLoggingInterceptor(log)),
	)

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       GRPCPort,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Any("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC started", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(slog.String("op", op))

	a.gRPCServer.GracefulStop()

	log.Info("gRPS server stopped", slog.Int("port", a.port))
}
