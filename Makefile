HTTP=pvz-service
GRPC=grpc-service
APP=$(HTTP) $(GRPC)
.PHONY: build run stop swag-gen unit-test integration-test load lint

build:
	docker-compose build

run:
	docker-compose up -d

stop:
	docker-compose down

swag-gen:
	swag init -g ../../cmd/http-server/main.go -o ./api -d ./internal/handlers

unit-test:
	go test ./... -coverprofile cover.out.tmp && \
    cat cover.out.tmp | grep -v "docs.go" |grep -v "router.go" | grep -v "mock_" | grep -v "grpc" | grep -v "db.go" | grep -v "config.go" | grep -v "logger.go" | grep -v "metrics.go" | grep -v "auth.go" | grep -v "dto" > cover.out && \
    rm cover.out.tmp && \
    go tool cover -func cover.out
	go tool cover -html=cover.out

integration-test:
	docker-compose up -d postgres
	go test -v ./tests/integration -tags=integration

load:
	k6 run tests/load/load.js

lint:
	golangci-lint run ./... --config=./.golangci.yaml
