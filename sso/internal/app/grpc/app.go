package grpc

import (
	"fmt"
	"log/slog"
	"net"

	authrpc "github.com/empfaze/golang_grpc/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

const OPERATION_TRACE_RUN = "grpcapp.Run"
const OPERATION_TRACE_STOP = "grpcapp.Stop"

func New(log *slog.Logger, authService authrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	logger := a.log.With(
		slog.String("op", OPERATION_TRACE_RUN),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", OPERATION_TRACE_RUN, err)
	}

	logger.Info("Starting grpc server...", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", OPERATION_TRACE_RUN, err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.With("op", OPERATION_TRACE_STOP).
		Info("Stopping grpc server...", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
