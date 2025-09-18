.PHONY: run build test docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

lint:
	golangci-lint run

mocks:
	mockgen -source=pkg/ads/client.go -destination=test/mocks/ads_client.go -package=mocks
	mockgen -source=pkg/crm/client.go -destination=test/mocks/crm_client.go -package=mocks