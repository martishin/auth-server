package suite

import (
	"context"
	"math/rand"
	"net"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	ssov1 "github.com/tty-monkey/auth-server-schemas/gen/go/sso"
	"github.com/tty-monkey/auth-server/internal/app"
	"github.com/tty-monkey/auth-server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost             = "localhost"
	grpcTimeout          = 5 * time.Second
	postgresStartTimeout = 5 * time.Second
	postgresOccurrence   = 2
	serviceStartWait     = 100 * time.Millisecond
	portOffset           = 20000
	portRandomPool       = 30000
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	postgresConnection := startPostgres(t)

	cfg := &config.Config{
		Env:      "local",
		PgConn:   postgresConnection,
		TokenTTL: 1 * time.Hour,
		GRPC: config.GRPCConfig{
			Port:    getRandomPort(),
			Timeout: grpcTimeout,
		},
	}

	go app.Start(cfg)
	time.Sleep(serviceStartWait)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

// Can start only one container at a time using testify suite
// https://testcontainers.com/guides/getting-started-with-testcontainers-for-go/#_reusing_the_containers_and_running_multiple_tests_as_a_suite
//
//nolint:lll // a url
func startPostgres(t *testing.T) string {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.1"),
		postgres.WithInitScripts(filepath.Join("testdata", "init-db.sql")),
		postgres.WithDatabase("test-auth"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(postgresOccurrence).WithStartupTimeout(postgresStartTimeout)),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable", "timezone=UTC")
	if err != nil {
		t.Fatalf("failed to get connection: %s", err)
	}

	return connStr
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

//nolint:gosec // okay to use in tests
func getRandomPort() int {
	return rand.Intn(portRandomPool) + portOffset
}
