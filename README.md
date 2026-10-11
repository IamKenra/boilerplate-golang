# svc-boilerplate-fiber-golang

Go Fiber service boilerplate with layered domain structure:

- `domain/<domain>/delivery/http`
- `domain/<domain>/usecase`
- `domain/<domain>/repository`
- shared modules: `utils/`, `entity/`, `valueobject/`

Default sample domain: `domain/masterdata`.

## Quick Start

```bash
cp .env.example .env
go run main.go
```

## Required Environment

Minimum variables for startup:

- `PORT`
- `DB_DIALECT` (`postgres` or `mysql`)
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USERNAME`

Recommended:

- `DB_PASSWORD`
- `APP_ENV` (`development` or `release`)
- `CORS_ALLOW_ORIGINS`
- `JWT_SECRET`
- `INTERNAL_API_KEY` (if using internal API key middleware)

Use `.env.example` as the source of truth for all available variables.

## Add Domain (Automation)

Create a new domain:

```bash
./add-domain.sh <domain_name>
```

Example:

```bash
./add-domain.sh auth
```

Generated files include:

- `domain/<domain>/usecase.go`
- `domain/<domain>/repository.go`
- `domain/<domain>/usecase/app.go`
- `domain/<domain>/usecase/api-health.go`
- `domain/<domain>/usecase/core-api.go`
- `domain/<domain>/usecase/processor.go`
- `domain/<domain>/repository/define.go`
- `domain/<domain>/repository/postgre-core.go`
- `domain/<domain>/repository/postgre-utils.go`
- `domain/<domain>/delivery/http/routes-http.go`
- `domain/<domain>/delivery/http/http-health.go`

The script also wires `main.go` with flow:

- `repo := NewPostgresRepository(postgreDB)`
- `uc := New<Domain>Usecase(repo)`
- `NewHttp<Domain>Handler(uc, app)`

## Health Check

```bash
curl --location 'http://localhost:3000/api/v1/masterdata/health'
```

## Database

Database access uses GORM with PostgreSQL and MySQL support. Set `DB_DIALECT=postgres` or `DB_DIALECT=mysql` in `.env`.

## Media Utilities

`utils/media` shells out to `cwebp` and `dwebp`; install those binaries before using WebP conversion.

## Verify Build

```bash
go test ./...
```
