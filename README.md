# boilerplate-golang

Go Fiber boilerplate with layered domain structure:

- `domain/<domain>/delivery/http`
- `domain/<domain>/usecase`
- `domain/<domain>/repository`
- shared modules: `utils/`, `entity/`, `valueobject/`

Default sample domain: `domain/masterdata`.

## Initial Setup
Set project name (also used as Go module/import root).
By default, script will ask interactively:

```bash
./setup-project.sh
```

You can still pass it directly:

```bash
./setup-project.sh --project svc-itasset-golang
./setup-project.sh backend_itasset
```

The script updates module/import references and tries `go mod tidy`.

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

## Run

```bash
export PORT=3000
go run main.go
```

## Health Check

Example for masterdata domain:

```bash
curl --location 'http://localhost:3000/api/v1/masterdata/health'
```

## Verify Build

```bash
go test ./...
```
