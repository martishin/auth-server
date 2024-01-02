package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storageConnection, migrationsPath, migrationTable string

	flag.StringVar(&storageConnection, "storage-connection", "", "storage connection string")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "name of migrations table")
	flag.Parse()

	if storageConnection == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	p := pgxMigrate.Postgres{}
	d, err := p.Open(fmt.Sprintf("%s&x-migrations-table=%s", storageConnection, migrationTable))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = d.Close(); err != nil {
			log.Fatalf("failed: close connection(migrations) due to error: %v", err)
		}
	}()

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, migrationTable, d)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	log.Println("migrations applied successfully")
}
