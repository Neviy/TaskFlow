run:
	go run ./cmd/api

test:
	go test ./...

build:
	go build ./cmd/api

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f api

docker-restart:
	docker compose down
	docker compose up -d --build

check:
	go fmt ./...
	go vet ./...
	go test ./...
	go build ./cmd/api