# Go Clean Architecture Template

Clean Architecture template for Golang

[![codecov](https://codecov.io/gh/DeSouzaRafael/go-clean-architecture-template/branch/main/graph/badge.svg?token=PNP58LYNPA)](https://codecov.io/gh/DeSouzaRafael/go-clean-architecture-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeSouzaRafael/go-clean-architecture-template)](https://goreportcard.com/report/github.com/DeSouzaRafael/go-clean-architecture-template)
[![License](https://img.shields.io/github/license/evrone/go-clean-template.svg)](https://github.com/DeSouzaRafael/go-clean-architecture-template/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/DeSouzaRafael/go-clean-architecture-template)](https://github.com/DeSouzaRafael/go-clean-architecture-template/releases/)

## Overview

Production-ready Go template that enforces Clean Architecture from the first line of code. The structure makes it impossible to accidentally couple business logic to frameworks or databases — the compiler will tell you if you do.

**What this gives you out of the box:**
- REST API with versioned routes (`/v0/`, `/v1/`, ...)
- PostgreSQL via GORM with a generic `BaseRepo[T]` for CRUD
- Structured JSON logging (zerolog)
- Request validation (go-playground/validator)
- Swagger UI at `/docs/*`
- Health check at `/health`
- Graceful shutdown on SIGTERM/SIGINT
- GitHub Actions CI: build, test, lint, security scan, tidy check
- Docker Compose for local development

## Table of Contents

- [Quick Start](#quick-start)
- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Dependency Injection](#dependency-injection)
- [Adding a New Domain](#adding-a-new-domain)
- [API Versioning](#api-versioning)
- [Development Commands](#development-commands)

## Quick Start

```sh
# 1. Copy env template
cp .env.example .env

# 2. Start postgres + app
make run
```

The API will be available at `http://localhost:8082` (or the `HTTP_PORT` you set in `.env`).

Swagger UI: `http://localhost:8082/docs/index.html`

## Environment Variables

Copy `.env.example` to `.env` and adjust as needed.

| Variable | Required | Default | Description |
|---|---|---|---|
| `ENV` | yes | — | `local` / `dev` / `prd` — controls AutoMigrate and CORS |
| `APP_NAME` | yes | — | Application name (used in logs) |
| `APP_VERSION` | yes | — | Application version |
| `HTTP_PORT` | yes | — | Port the server listens on |
| `LOG_LEVEL` | yes | — | `debug` / `info` / `warn` / `error` |
| `PG_URL` | yes | — | Full PostgreSQL DSN |
| `PG_MAX_OPEN_CONNS` | no | `10` | Connection pool max open |
| `PG_MAX_IDLE_CONNS` | no | `5` | Connection pool max idle |
| `PG_CONN_MAX_LIFETIME_SEC` | no | `3600` | Connection max lifetime in seconds |

> When `ENV=prd`, AutoMigrate is skipped at startup. Use [golang-migrate](https://github.com/golang-migrate/migrate) or similar for production schema changes.

## Project Structure

```
.
├── cmd/app/                        # Binary entry point
├── config/                         # Config struct loaded from env
├── docs/                           # Auto-generated Swagger (do not edit)
├── infra/
│   ├── httpserver/                 # HTTP server with graceful shutdown
│   ├── logger/                     # zerolog wrapper
│   ├── postgres/
│   │   ├── model/                  # GORM models + entity↔model mappers
│   │   └── repository/             # BaseRepo[T] + domain-specific repos
│   └── validator/                  # Request validation wrapper
├── internal/
│   ├── app/                        # Wiring: infra + repos + use cases + server
│   ├── controller/rest/
│   │   ├── input/                  # Request DTOs
│   │   ├── output/                 # Response DTOs
│   │   └── routers/v0/             # Route handlers
│   ├── entity/                     # Pure domain structs (zero framework tags)
│   └── usecase/                    # Business logic + interface contracts
└── mocks/                          # Generated mocks (go.uber.org/mock)
```

### Key Design Decisions

**`internal/entity/` — pure domain structs**

Entities have no GORM tags, no JSON tags, no `TableName()`. They are plain Go structs that carry domain data between layers. Serialization is handled by output DTOs; persistence is handled by GORM models.

```go
type UserEntity struct {
    ID        uuid.UUID
    Name      string
    Phone     string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

**`infra/postgres/model/` — GORM models**

GORM concerns live here: column tags, `TableName()`, soft-delete type, default expressions. Each model file also contains `ToXxxModel` and `ToXxxEntity` mapper functions.

```go
type UserModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name      string
    Phone     string
    CreatedAt time.Time      `gorm:"<-:create"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**`infra/postgres/repository/BaseRepo[T]` — generic CRUD**

`BaseRepo[T any]` provides `Get`, `Create`, `Update`, `Delete` for any GORM model. Domain repos embed it and expose only domain-typed methods:

```go
type UserRepo struct {
    *BaseRepo[model.UserModel]
}

func (r *UserRepo) GetById(ctx context.Context, e entity.UserEntity) (entity.UserEntity, error) {
    m, err := r.BaseRepo.Get(ctx, e.ID)
    return model.ToUserEntity(m), err
}
```

**`internal/usecase/interfaces.go` — single source of truth**

All cross-layer contracts live here. The `UseCases` aggregator is the only type injected into the router — controllers never depend on concrete implementations.

```go
type User interface {
    CreateUser(context.Context, entity.UserEntity) (entity.UserEntity, error)
    UpdateUser(context.Context, entity.UserEntity) error
    DeleteUser(context.Context, entity.UserEntity) error
    GetUserById(context.Context, entity.UserEntity) (entity.UserEntity, error)
}

type UserRepo interface {
    Create(context.Context, entity.UserEntity) (entity.UserEntity, error)
    Update(context.Context, entity.UserEntity) error
    DeleteById(context.Context, entity.UserEntity) error
    GetById(context.Context, entity.UserEntity) (entity.UserEntity, error)
}

type UseCases interface {
    UserUseCase() User
}
```

## Architecture

The dependency rule is enforced by Go's package system. Each layer imports only inward:

```
cmd/app
  └── internal/app          (composition root — imports everything once)
        ├── internal/usecase      (business logic)
        │     └── internal/entity (domain structs)
        └── infra/                (postgres, logger, validator, httpserver)
              └── internal/entity (for model mappers)

internal/controller/rest
  └── internal/usecase interfaces (never concrete types)
```

`internal/` packages cannot be imported by external modules (Go enforces this). The only place that crosses all layers is `internal/app/app.go` — the composition root.

## Dependency Injection

Dependencies flow inward through constructors. Use cases receive repository interfaces, not concrete implementations:

```go
// Declare what you need
type UserRepo interface {
    GetById(ctx context.Context, user entity.UserEntity) (entity.UserEntity, error)
}

// Receive it via constructor
type UserUseCase struct {
    repo UserRepo
}

func NewUser(r UserRepo) *UserUseCase {
    return &UserUseCase{repo: r}
}
```

To swap PostgreSQL for another database, write a new struct that satisfies `UserRepo` and wire it in `internal/app/app.go`. No use-case code changes.

## Adding a New Domain

Example: adding `Product`.

**1. Domain entity** — `internal/entity/product_entity.go`
```go
type ProductEntity struct {
    ID    uuid.UUID
    Name  string
    Price float64
}
```

**2. GORM model + mappers** — `infra/postgres/model/product_model.go`
```go
type ProductModel struct {
    ID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name  string
    Price float64
}
func (ProductModel) TableName() string { return "product" }
func ToProductModel(e entity.ProductEntity) ProductModel { ... }
func ToProductEntity(m ProductModel) entity.ProductEntity { ... }
```

**3. Repository interface** — add to `internal/usecase/interfaces.go`
```go
ProductRepo interface {
    Create(context.Context, entity.ProductEntity) (entity.ProductEntity, error)
    GetById(context.Context, entity.ProductEntity) (entity.ProductEntity, error)
}
```

**4. Repository implementation** — `infra/postgres/repository/product_repository.go`
```go
type ProductRepo struct{ *BaseRepo[model.ProductModel] }
```

**5. Use case interface + implementation** — `internal/usecase/interfaces.go` + `internal/usecase/product_usecase.go`

**6. Wire** — `internal/app/app.go`: add to `AutoMigrate` and `NewAppUseCases`

**7. Handlers + DTOs** — `internal/controller/rest/routers/v0/product_view.go`, `input/`, `output/`

**8. Regenerate mocks**
```sh
make mock
```

## API Versioning

Routes are namespaced by version. To add v1:

1. Create `internal/controller/rest/routers/v1/` with new handlers
2. Register in `internal/controller/rest/router.go`:

```go
v0.NewUserRoutes(h, l, v, uc.UserUseCase())
v1.NewProductRoutes(h, l, v, uc.ProductUseCase())
```

v0 and v1 can coexist indefinitely.

## Development Commands

```sh
make run          # start postgres (docker) + run app locally
make test         # go test -v -cover -race ./internal/...
make swag         # regenerate Swagger docs
make mock         # regenerate mocks from interfaces.go
make linter-golangci  # run golangci-lint

# Docker
docker-compose up -d          # start postgres + app
docker-compose up -d postgres # start only postgres
```
