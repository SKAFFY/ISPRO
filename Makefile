.PHONY: start test bench clean build-server build-regexp lint coverage migrate-up migrate-down docker-up docker-down docker-build generate-api di-generate swagger

BINARY_NAME=ispro-app
POSTGRES_DSN=postgres://ispro:ispro@localhost:5432/ispro?sslmode=disable

start: build-server
	./bin/$(BINARY_NAME)

build-server:
	go build -o bin/$(BINARY_NAME) ./cmd/server

build-regexp:
	go build -o bin/regexp ./cmd/regexp

run-regexp: build-regexp
	./bin/regexp

test:
	go test ./...

bench:
	go test -bench=. ./...

migrate-up:
	goose -dir=internal/migrations postgres $(POSTGRES_DSN) up

migrate-down:
	goose -dir=internal/migrations postgres $(POSTGRES_DSN) down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker build -t $(BINARY_NAME) .

generate-api:
	rm -rf cmd/ispro-app-server
	mkdir -p cmd/ispro-app-server
	cd cmd && swagger generate server -A ispro-app -f ../api/openapi.yaml --principal string --target ispro-app-server
	rm -rf cmd/ispro-app-server/cmd
	rm -rf cmd/ispro-app-server/internal
	cp -r cmd/ispro-app-server/models/* internal/models/ 2>/dev/null || true
	cp -r cmd/ispro-app-server/restapi/* internal/restapi/ 2>/dev/null || true
	rm -rf cmd/ispro-app-server

di-generate:
	digen generate

swagger:
	swagger serve api/openapi.yaml -p 8081