.PHONY: run build swagger migrate

## run: start the API server on the default port
run:
	go run ./cmd apiserver --port 8080

## build: compile the binary
build:
	go build -o bin/arena ./cmd

## swagger: regenerate Swagger docs (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	swag init -g cmd/main.go -o docs

## migrate: run database migrations (DATABASE_URL env var must be set)
migrate:
	go run ./cmd migrate --database-url "$(DATABASE_URL)"
