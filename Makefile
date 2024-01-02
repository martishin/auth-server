run-postgres:
	docker-compose up db

test:
	go test ./...

migrate:
	go run ./cmd/migrator --storage-connection="postgres://postgres:postgres@localhost:5432/auth?sslmode=disable&timezone=UTC&connect_timeout=5" --migrations-path=./migrations
