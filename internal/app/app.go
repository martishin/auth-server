package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/tty-monkey/auth-server/internal/app/grpc"
	"github.com/tty-monkey/auth-server/internal/services/auth"
)

type App struct {
	GRPCserver *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	// storagePath string,
	_ string,
	tokenTTL time.Duration,
) *App {
	grpcAuth := auth.New(log, nil, nil, nil, tokenTTL)
	grpcApp := grpcapp.New(log, grpcAuth, grpcPort)

	return &App{
		GRPCserver: grpcApp,
	}
}
