package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"vizapSSO/internal/app"
	"vizapSSO/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting SSO app",
		slog.String("env", cfg.Env),
		slog.Int("port", cfg.GRPC.Port),
	)

	application := app.New(log, cfg.GRPC.Port, cfg.AccessTokenTTL, cfg.RefreshTokenTTL, cfg.Postgres)

	go application.GRPSServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("stopping SSO app")

	application.GRPSServer.Stop()

	log.Info("SSO app stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
