.PHONY: run test bench clean

run:
	go run ./cmd/regexp

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./cmd/regexp/...

bench:
	go test -bench=. -benchmem ./cmd/regexp/...

lint:
	golangci-lint run ./...

coverage:
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f coverage.out coverage.html

build:
	go build -o bin/regexp ./cmd/regexp