package app

import (
	"log/slog"
	"strconv"
	"time"
	grpcapp "vizapSSO/internal/app/grpc"
	"vizapSSO/internal/config"
	"vizapSSO/internal/services/auth"
	"vizapSSO/internal/storage/postgres"
)

type App struct {
	GRPSServer *grpcapp.App
}

func New(log *slog.Logger, GRPCPort int, AccessTokenTTL time.Duration, RefreshTokenTTL time.Duration, postgresConfig config.PostgresConfig) *App {
	storage, err := postgres.New(postgresConfig.Host, strconv.Itoa(postgresConfig.Port), postgresConfig.User, postgresConfig.Password, postgresConfig.DBName)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, storage, storage, AccessTokenTTL, RefreshTokenTTL)

	grpcApp := grpcapp.New(log, authService, GRPCPort)

	return &App{
		GRPSServer: grpcApp,
	}
}
