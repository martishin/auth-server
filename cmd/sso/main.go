package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tty-monkey/auth-server/internal/app"
	"github.com/tty-monkey/auth-server/internal/config"
	"github.com/tty-monkey/auth-server/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("starting application", slog.Any("config", cfg))

	application := app.New(logger, cfg.GRPC.Port, cfg.PostgresConnection, cfg.TokenTTL)

	go application.GRPCserver.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	logger.Info("received signal", slog.String("signal", sig.String()))

	application.GRPCserver.Stop()

	logger.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = setupPrettySlog()
	case envDev:
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
