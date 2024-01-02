package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tty-monkey/auth-server/internal/domain/models"
	"github.com/tty-monkey/auth-server/internal/storage"
)

type Storage struct {
	DB *pgxpool.Pool
}

// New creates new instance of the PostgreSQL Storage.
func New(storageConnection string) (*Storage, error) {
	const op = "storage.postgresql.New"

	pool, err := pgxpool.New(context.Background(), storageConnection)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: pool}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	const op = "storage.postgresql.SaveUser"
	// new ctx with timeout can be created
	// ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	// defer cancel()

	stmt := `
		INSERT INTO users (email, pass_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var userID int64
	err := s.DB.QueryRow(ctx, stmt, email, passwordHash).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgresql.User"

	query := `
		SELECT id, email, pass_hash
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := s.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgresql.IsAdmin"

	query := `
		SELECT is_admin
		FROM users
		WHERE id = $1
	`

	var isAdmin bool
	err := s.DB.QueryRow(ctx, query, userID).Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgresql.App"

	query := `
		SELECT id, name, secret
		FROM apps
		WHERE id = $1
	`

	var app models.App
	err := s.DB.QueryRow(ctx, query, appID).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
