package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/tty-monkey/auth-server/internal/app/grpc"
)

type App struct {
	GRPCserver *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	// storagePath string,
	_ string,
	// tokenTTL time.Duration,
	_ time.Duration,
) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCserver: grpcApp,
	}
}
