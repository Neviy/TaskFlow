APP_NAME=taskflow
API_PATH=./cmd/api

.PHONY: run test build fmt vet tidy check clean

run:
	go run $(API_PATH)

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) $(API_PATH)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

check:
	go fmt ./...
	go vet ./...
	go test ./...
	go build $(API_PATH)

clean:
	go clean