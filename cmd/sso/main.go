package main

import (
	"io"
	stlog "log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	application := app.New(log, cfg.GRPC.Port, cfg.AccessTokenTTL, cfg.RefreshTokenTTL, cfg.Postgres)

	go application.GRPSServer.MustRun()

	Print(cfg, log)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("stopping SSO app")

	application.GRPSServer.Stop()

	log.Info("SSO app stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func Print(cfg *config.Config, log *slog.Logger) {
	time.Sleep(time.Second * 3)

	stlog.Println("==========SSO APP STARTED===========")
	stlog.Printf("|gRPC PORT................%d\n", cfg.GRPC.Port)
	stlog.Printf("|POSTGRESQL HOST..........%s\n", cfg.Postgres.Host)
	stlog.Printf("|POSTGRESQL PORT..........%d\n", cfg.Postgres.Port)
	stlog.Printf("|ACCESS TOKEN TTL.........%s\n", cfg.AccessTokenTTL)
	stlog.Printf("|REFRESH TOKEN TTL........%s\n", cfg.RefreshTokenTTL)
	stlog.Printf("|ENV CONFIG...............%s\n", cfg.Env)
	stlog.Println("====================================")

}
