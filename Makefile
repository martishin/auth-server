run-postgres:
	docker-compose up db

test:
	go test ./...
