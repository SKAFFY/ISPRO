.PHONY: run test bench clean build lint coverage migrate-up migrate-down docker-up docker-down

BINARY_NAME=ispro-app
POSTGRES_DSN=postgres://ispro:ispro@localhost:5432/ispro?sslmode=disable

run:
	go run ./cmd/server

run-regexp:
	go run ./cmd/regexp

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run ./...

coverage:
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f coverage.out coverage.html bin/$(BINARY_NAME)

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

build-regexp:
	go build -o bin/regexp ./cmd/regexp

migrate-up:
	goose -dir=migrations postgres $(POSTGRES_DSN) up

migrate-down:
	goose -dir=migrations postgres $(POSTGRES_DSN) down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker build -t $(BINARY_NAME) .

generate:
	swagger generate server -A ispro-app -f api/openapi.yaml --principal string