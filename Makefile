run-postgres:
ifeq ($(DETACHED),true)
	docker-compose up db -d
else
	docker-compose up db
endif

run-app:
ifeq ($(DETACHED),true)
	docker-compose up app -d
else
	docker-compose up app
endif

test:
	go test ./...

migrate:
	go run ./cmd/migrator --storage-connection="postgres://postgres:postgres@localhost:5432/auth?sslmode=disable&timezone=UTC&connect_timeout=5" --migrations-path=./migrations

dump-postgres:
	@export PGPASSWORD="postgres" && \
	pg_dump -h localhost -U postgres --schema-only -d auth > schema_dump.sql

build:
	docker build -t auth-server .

run:
	export CONFIG_PATH=./configs/local.yaml; \
	$(MAKE) run-postgres DETACHED=true; \
	$(MAKE) migrate; \
	$(MAKE) run-app DETACHED=true

stop:
	docker-compose stop

down:
	docker-compose down

remove:
	docker-compose down -v

logs:
	docker-compose logs -f
