# Go Clean Architecture Template

Clean Architecture template for Golang

[![CI](https://github.com/DeSouzaRafael/go-clean-architecture-template/actions/workflows/ci.yml/badge.svg)](https://github.com/DeSouzaRafael/go-clean-architecture-template/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/DeSouzaRafael/go-clean-architecture-template/branch/main/graph/badge.svg?token=PNP58LYNPA)](https://codecov.io/gh/DeSouzaRafael/go-clean-architecture-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeSouzaRafael/go-clean-architecture-template)](https://goreportcard.com/report/github.com/DeSouzaRafael/go-clean-architecture-template)
[![Go Version](https://img.shields.io/github/go-mod/go-version/DeSouzaRafael/go-clean-architecture-template)](go.mod)
[![License](https://img.shields.io/github/license/DeSouzaRafael/go-clean-architecture-template)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/DeSouzaRafael/go-clean-architecture-template)](https://github.com/DeSouzaRafael/go-clean-architecture-template/releases/)

## Overview

Production-ready Go template that enforces Clean Architecture from the first line of code. The structure makes it impossible to accidentally couple business logic to frameworks or databases — the compiler will tell you if you do.

**What this gives you out of the box:**

| Feature | Details |
|---|---|
| REST API | Versioned routes (`/v0/`, `/v1/`, ...) via [Echo](https://echo.labstack.com/) |
| Database | PostgreSQL via GORM with generic `BaseRepo[T]` |
| Logging | Structured JSON via [zerolog](https://github.com/rs/zerolog) |
| Validation | Request validation via [go-playground/validator](https://github.com/go-playground/validator) |
| API Docs | Swagger UI at `/docs/index.html` |
| Observability | Health check at `/health` |
| Shutdown | Graceful on SIGTERM/SIGINT |
| CI | Build, test, lint, security scan, tidy check |
| Testing | 80%+ coverage with unit tests, mocks via [go.uber.org/mock](https://github.com/uber-go/mock) |

## Table of Contents

- [Quick Start](#quick-start)
- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)
- [Request Flow](#request-flow)
- [Architecture](#architecture)
- [Dependency Injection](#dependency-injection)
- [Testing](#testing)
- [Adding a New Domain](#adding-a-new-domain)
- [API Versioning](#api-versioning)
- [Development Commands](#development-commands)

## Quick Start

### Option 1 — Docker (recommended)

```sh
cp .env.example .env
docker-compose up -d
```

App starts at `http://localhost:8082`. Postgres on port `5434`.

### Option 2 — Local (requires Go 1.25+)

```sh
cp .env.example .env
# Start only postgres
docker-compose up -d postgres
# Run app
make run
```

Swagger UI: `http://localhost:8082/docs/index.html`

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `ENV` | yes | — | `local` / `dev` / `prd` — controls AutoMigrate and CORS |
| `APP_NAME` | yes | — | Application name |
| `APP_VERSION` | yes | — | Application version |
| `HTTP_PORT` | yes | `8082` | Port the server binds to |
| `LOG_LEVEL` | yes | — | `debug` / `info` / `warn` / `error` |
| `PG_URL` | yes | — | PostgreSQL DSN (e.g. `postgres://user:pass@host:5432/db`) |
| `PG_MAX_OPEN_CONNS` | no | `10` | Connection pool max open |
| `PG_MAX_IDLE_CONNS` | no | `5` | Connection pool max idle |
| `PG_CONN_MAX_LIFETIME_SEC` | no | `3600` | Connection lifetime in seconds |

> `ENV=prd` skips `AutoMigrate` at startup. Use [golang-migrate](https://github.com/golang-migrate/migrate) for production schema management.

## Project Structure

```
.
├── cmd/app/                        # Binary entry point (main.go)
├── config/                         # Typed config loaded from env vars
├── docs/                           # Auto-generated Swagger — do not edit
├── infra/
│   ├── httpserver/                 # net/http wrapper with graceful shutdown
│   ├── logger/                     # zerolog implementation of logger.Interface
│   ├── postgres/
│   │   ├── model/                  # GORM models + bidirectional mappers
│   │   └── repository/             # Generic BaseRepo[T] + domain repos
│   └── validator/                  # go-playground/validator wrapper
├── internal/
│   ├── app/                        # Composition root — wires all layers
│   ├── controller/rest/
│   │   ├── input/                  # Request DTOs with validation
│   │   ├── output/                 # Response DTOs and error helpers
│   │   └── routers/v0/             # Versioned route handlers
│   ├── entity/                     # Pure domain structs (no framework tags)
│   └── usecase/                    # Business logic + interface contracts
└── mocks/                          # go.uber.org/mock generated mocks
```

### Key Packages

**`internal/entity/`** — Pure domain structs with no GORM tags, no JSON tags. Used across all layers as the canonical data type. Serialization is handled by output DTOs; persistence by GORM models.

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

**`infra/postgres/model/`** — GORM concerns are isolated here. Each model includes column tags, `TableName()`, soft-delete via `gorm.DeletedAt`, and bidirectional mapper functions.

```go
type UserModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name      string
    Phone     string
    CreatedAt time.Time      `gorm:"<-:create"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string { return "user" }
func ToUserModel(e entity.UserEntity) UserModel   { ... }
func ToUserEntity(m UserModel) entity.UserEntity  { ... }
```

**`infra/postgres/repository/BaseRepo[T]`** — Generic CRUD. Domain repos embed it and add entity-typed methods:

```go
type UserRepo struct{ *BaseRepo[model.UserModel] }

func (r *UserRepo) GetById(ctx context.Context, e entity.UserEntity) (entity.UserEntity, error) {
    m, err := r.BaseRepo.Get(ctx, e.ID)
    return model.ToUserEntity(m), err
}
```

**`internal/usecase/interfaces.go`** — Single source of truth for all cross-layer contracts. Controllers and use cases depend only on these interfaces.

## Request Flow

```
HTTP Request
    │
    ▼
Echo router (internal/controller/rest/router.go)
    │  middleware: CORS, Recover
    ▼
Route handler (internal/controller/rest/routers/v0/user_view.go)
    │  1. Bind JSON → input DTO
    │  2. Validate (infra/validator)
    │  3. Map input → entity.UserEntity
    ▼
Use case (internal/usecase/user_usecase.go)
    │  Business rules, orchestration
    ▼
Repository interface (internal/usecase/interfaces.go)
    ▼
Repository implementation (infra/postgres/repository/user_repository.go)
    │  entity → model (mapper)
    │  GORM query → PostgreSQL
    │  model → entity (mapper)
    ▼
Back up through use case → handler
    │  Map entity → output DTO
    ▼
JSON Response
```

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                  internal/entity                    │  no deps
├─────────────────────────────────────────────────────┤
│                  internal/usecase                   │  → entity
├─────────────────────────────────────────────────────┤
│              internal/controller/rest               │  → usecase interfaces
├─────────────────────────────────────────────────────┤
│  infra/ (postgres, logger, validator, httpserver)   │  → entity (models)
├─────────────────────────────────────────────────────┤
│                 internal/app                        │  imports all (once)
└─────────────────────────────────────────────────────┘
```

Each layer imports only inward. `internal/` is Go-enforced: external packages cannot import it. The only place that crosses all layers is `internal/app/app.go`.

## Dependency Injection

Constructors receive dependencies as interfaces. Concrete types never leak upward:

```go
type UserRepo interface {
    GetById(ctx context.Context, user entity.UserEntity) (entity.UserEntity, error)
}

type UserUseCase struct{ repo UserRepo }

func NewUser(r UserRepo) *UserUseCase {
    return &UserUseCase{repo: r}
}
```

To replace PostgreSQL: implement `UserRepo` with any backend, wire it in `internal/app/app.go`. Zero changes in `internal/usecase`.

## Testing

Tests are co-located with the source they cover. Mocks are generated from `internal/usecase/interfaces.go`:

```sh
make mock   # regenerate after changing interfaces
make test   # go test -v -cover -race ./internal/...
```

**Coverage targets:**

| Package | Strategy |
|---|---|
| `internal/usecase` | Unit tests with mock repository |
| `internal/controller/rest/routers/v0` | Handler tests with mock use case |
| `config`, `infra/validator`, `infra/logger` | Unit tests |
| `infra/postgres/model` | Mapper round-trip tests |
| `infra/httpserver` | Server option tests |
| `infra/postgres/repository` | Integration tests (requires DB) |

Packages excluded from coverage: `mocks/`, `docs/`, `cmd/`, `infra/postgres/repository/`, `internal/app/`.

## Adding a New Domain

Example: `Product`.

**1.** `internal/entity/product_entity.go`
```go
type ProductEntity struct {
    ID    uuid.UUID
    Name  string
    Price float64
}
```

**2.** `infra/postgres/model/product_model.go`
```go
type ProductModel struct {
    ID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name  string
    Price float64
}
func (ProductModel) TableName() string { return "product" }
func ToProductModel(e entity.ProductEntity) ProductModel  { ... }
func ToProductEntity(m ProductModel) entity.ProductEntity { ... }
```

**3.** Add to `internal/usecase/interfaces.go`:
```go
Product interface {
    CreateProduct(context.Context, entity.ProductEntity) (entity.ProductEntity, error)
    GetProductById(context.Context, entity.ProductEntity) (entity.ProductEntity, error)
}

ProductRepo interface {
    Create(context.Context, entity.ProductEntity) (entity.ProductEntity, error)
    GetById(context.Context, entity.ProductEntity) (entity.ProductEntity, error)
}

// Add to UseCases aggregator:
ProductUseCase() Product
```

**4.** `infra/postgres/repository/product_repository.go`
```go
type ProductRepo struct{ *BaseRepo[model.ProductModel] }
```

**5.** `internal/usecase/product_usecase.go` — business logic

**6.** `internal/app/app.go` — register in `AutoMigrate` + `NewAppUseCases`

**7.** `internal/controller/rest/routers/v0/product_view.go` — handlers

**8.** `internal/controller/rest/input/product_input.go`, `output/product_output.go` — DTOs

**9.**
```sh
make mock  # regenerate mocks
```

## API Versioning

Routes namespace by version prefix. To add v1:

```sh
mkdir -p internal/controller/rest/routers/v1
```

Register in `internal/controller/rest/router.go`:

```go
v0.NewUserRoutes(h, l, v, uc.UserUseCase())
v1.NewProductRoutes(h, l, v, uc.ProductUseCase())
```

Both versions coexist. v0 routes remain unchanged.

## Development Commands

```sh
# Run
make run              # start postgres + run app

# Test
make test             # run tests with coverage and race detector

# Docs
make swag             # regenerate Swagger from annotations

# Mocks
make mock             # regenerate mocks from interfaces.go

# Lint
make linter-golangci  # run golangci-lint

# Install dev tools
make bin-dependencies # install swag + mockgen binaries

# Docker
docker-compose up -d              # start postgres + app
docker-compose up -d postgres     # postgres only
docker-compose down               # stop all
docker-compose logs -f app        # stream app logs
```
