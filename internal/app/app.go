package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/tty-monkey/auth-server/internal/app/grpc"
	"github.com/tty-monkey/auth-server/internal/services/auth"
	"github.com/tty-monkey/auth-server/internal/storage/postgresql"
)

type App struct {
	GRPCserver *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	storageConnection string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgresql.New(storageConnection)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, port)

	return &App{
		GRPCserver: grpcApp,
	}
}
