package main

import (
	"github.com/empfaze/golang_grpc/sso/internal/config"
	"github.com/empfaze/golang_grpc/sso/internal/logger"
)

func main() {
	config := config.MustLoad()
	logger := logger.SetupLogger(config.Env)

	logger.Info("Starting server...")
}
