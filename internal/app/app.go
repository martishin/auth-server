package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcapp "github.com/tty-monkey/auth-server/internal/app/grpc"
	"github.com/tty-monkey/auth-server/internal/config"
	"github.com/tty-monkey/auth-server/internal/lib/logger/handlers/slogpretty"
	"github.com/tty-monkey/auth-server/internal/services/auth"
	"github.com/tty-monkey/auth-server/internal/storage/postgresql"
)

type App struct {
	GRPCserver *grpcapp.App
	Storage    *postgresql.Storage
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Start(cfg *config.Config) {
	logger := setupLogger(cfg.Env)

	logger.Info("starting application", slog.Any("config", cfg))

	application := New(logger, cfg.GRPC.Port, cfg.PgConn, cfg.TokenTTL)

	go application.GRPCserver.MustRun()

	defer func() {
		logger.Info("closing database connection")
		application.Storage.DB.Close()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	logger.Info("received signal", slog.String("signal", sig.String()))

	application.GRPCserver.Stop()

	logger.Info("application stopped")
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
		Storage:    storage,
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = setupPrettySlog()
	case envDev, envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlopOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
