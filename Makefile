include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

.PHONY: help

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres
.PHONY: compose-up

swag-v0: ### swag init
	swag init -g internal/controller/api/router.go 
.PHONY: swag-v0

run: swag-v0 ### swag run
	go mod tidy && go mod download && CGO_ENABLED=0 go run ./cmd/app 
.PHONY: run

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock: ### run mockgen
	mockgen -source ./internal/usecase/interfaces.go -package mocks > ./mocks/mock_user_usecase.go
.PHONY: mock

bin-dependencies:
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest
.PHONY: bin-dependencies

#mockgen -source ./infra/postgres/postgres.go -package infra_test > ./mocks/mock_postgres.go