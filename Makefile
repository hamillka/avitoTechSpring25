HTTP=pvz-service
GRPC=grpc-service
APP=$(HTTP) $(GRPC)
.PHONY: build run stop swag-gen unit-test integration-test

build:
	docker-compose build $(APP)

run:
	docker-compose up -d $(APP)

stop:
	docker-compose down

swag-gen:
	swag init -g ../../cmd/http-server/main.go -o ./api -d ./internal/handlers

unit-test:
	go test ./tests/units/... -coverprofile=coverage.out -coverpkg=./...
	go tool cover -html=coverage.out

integration-test:
	docker-compose up -d postgres
	go test -v ./tests/integration/...
