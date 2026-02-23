# boilerplate-golang

Minimal clean boilerplate that mirrors the original architecture:

- `domain/` (delivery/usecase/repository)
- `utils/`
- `entity/`, `valueobject/`

Currently includes only `domain/masterdata` as an example.
All handlers return simple health responses so the project compiles and runs.

## Initial Setup
Set project name first (it will also be used as Go module/import root):
```
./setup-project.sh --project svc-itasset-golang
./setup-project.sh backend_itasset
```
This command also runs `go mod tidy` to install/update required dependencies.

## Add Domains (Automation)
Use the script in this folder to add new domains into this boilerplate:
```
./add-domain.sh <domain_name>
```

Example:
```
./add-domain.sh auth
```

## Run
```
export PORT=3000
go run main.go
```
