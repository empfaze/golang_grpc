package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/empfaze/golang_grpc/sso/internal/app"
	"github.com/empfaze/golang_grpc/sso/internal/config"
	"github.com/empfaze/golang_grpc/sso/internal/logger"
)

func main() {
	config := config.MustLoad()
	logger := logger.SetupLogger(config.Env)

	logger.Info("Starting server...")

	application := app.New(logger, config.GRPC.Port, config.StoragePath, config.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	logger.Info("Application stopping...", slog.String("signal", sig.String()))

	application.GRPCSrv.Stop()

	logger.Info("Application has been gracefully stopped")
}
