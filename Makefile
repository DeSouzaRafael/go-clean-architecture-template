include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

.PHONY: help

compose-up: ### Run app in docker
	docker-compose up --build -d postgres app && \
	docker-compose logs -f app
.PHONY: compose-up

swag: ### swag init
	swag init -g internal/controller/rest/router.go
.PHONY: swag

run: swag ### run project
	docker-compose up --build -d postgres && \
	go mod tidy && go mod download && \
	CGO_ENABLED=0 go run ./cmd/app
.PHONY: run

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

mock: ### run mockgen
	mockgen -source ./internal/usecase/interfaces.go -package mocks > ./mocks/mock_user_usecase.go
.PHONY: mock

bin-dependencies: ### install dev dependencies
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest
.PHONY: bin-dependencies

golangci-lint-install:
	@echo "Installing golangci-lint v1.59.1..."
	@if [ ! -f "$(LOCAL_BIN)/golangci-lint" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) v1.59.1; \
	fi
.PHONY: golangci-lint-install
