#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./add-domain.sh <domain_name>

Example:
  ./add-domain.sh auth
EOF
}

if [ $# -lt 1 ]; then
  usage
  exit 1
fi

DOMAIN_RAW="$1"
DOMAIN="$(echo "${DOMAIN_RAW}" | tr '[:upper:]' '[:lower:]')"

if ! [[ "${DOMAIN}" =~ ^[a-z][a-z0-9_]*$ ]]; then
  echo "invalid domain name: ${DOMAIN}" >&2
  exit 1
fi

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOILER_DIR="${ROOT_DIR}"

GO_MOD_FILE="${BOILER_DIR}/go.mod"
if [ ! -f "${GO_MOD_FILE}" ]; then
  echo "go.mod not found in boilerplate" >&2
  exit 1
fi

MODULE_PATH="$(awk '$1=="module"{print $2; exit}' "${GO_MOD_FILE}")"
if [ -z "${MODULE_PATH}" ]; then
  echo "failed to read module path from boilerplate/go.mod" >&2
  exit 1
fi

DOMAIN_DIR="${BOILER_DIR}/domain/${DOMAIN}"
if [ -d "${DOMAIN_DIR}" ]; then
  echo "domain already exists: ${DOMAIN_DIR}" >&2
  exit 1
fi

TITLE_DOMAIN="$(echo "${DOMAIN}" | awk '{print toupper(substr($0,1,1)) substr($0,2)}')"

write_file() {
  local path="$1"
  shift
  mkdir -p "$(dirname "${path}")"
  cat > "${path}" <<EOF
$*
EOF
}

write_file "${DOMAIN_DIR}/usecase.go" "package ${DOMAIN}

type Usecase interface {
	Health() (string, error)
}"

write_file "${DOMAIN_DIR}/usecase/app.go" "package usecase

import \"${MODULE_PATH}/domain/${DOMAIN}\"

/*
*
the struct of usecases
*/
type ${DOMAIN}Usecase struct {
	postgreRepository ${DOMAIN}.PostgreRepository
}

/*
*
the initiator function for usecase
*/
func New${TITLE_DOMAIN}Usecase(postgres ${DOMAIN}.PostgreRepository) ${DOMAIN}.Usecase {
	return &${DOMAIN}Usecase{
		postgreRepository: postgres,
	}
}"

write_file "${DOMAIN_DIR}/repository.go" "package ${DOMAIN}

type PostgreRepository interface {
	Health() error
	GenerateID() (uint64, error)
	GenerateUUID() (string, error)
}"

write_file "${DOMAIN_DIR}/usecase/api-health.go" "package usecase

func (${DOMAIN} ${DOMAIN}Usecase) Health() (string, error) {
	err := ${DOMAIN}.postgreRepository.Health()
	if err != nil {
		return \"\", err
	}

	return \"Service is healthy\", nil
}"

write_file "${DOMAIN_DIR}/usecase/core-api.go" "package usecase

// Place business API orchestration methods in this file."

write_file "${DOMAIN_DIR}/usecase/processor.go" "package usecase

// Place query builder processors in this file."

write_file "${DOMAIN_DIR}/repository/define.go" "package repository

import (
	\"database/sql\"

	\"${MODULE_PATH}/domain/${DOMAIN}\"
)

const POSTGRE = \"postgre\"

const SELECT_DISTINCT = \"select-distinct\"
const SELECT = \"select\"
const INSERT = \"insert\"
const UPDATE = \"update\"
const DELETE = \"delete\"

type postgres${TITLE_DOMAIN}Repository struct {
	sqlDB *sql.DB
}

func NewPostgresRepository(db *sql.DB) ${DOMAIN}.PostgreRepository {
	return &postgres${TITLE_DOMAIN}Repository{db}
}"

write_file "${DOMAIN_DIR}/repository/postgre-core.go" "package repository

// Place domain-specific PostgreSQL query methods in this file."

write_file "${DOMAIN_DIR}/repository/postgre-utils.go" "package repository

import (
	\"log\"
	\"os\"
	\"strconv\"

	\"github.com/bwmarrin/snowflake\"
	\"github.com/google/uuid\"
)

var node *snowflake.Node

func init() {
	var err error
	nodeID := int64(1)
	if raw := os.Getenv(\"NODE_ID\"); raw != \"\" {
		nodeID, err = strconv.ParseInt(raw, 10, 64)
		if err != nil || nodeID < 0 || nodeID > 1023 {
			log.Fatalf(\"invalid NODE_ID: %q\", raw)
		}
	}
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf(\"Failed to create snowflake node: %v\", err)
	}
}

func (db *postgres${TITLE_DOMAIN}Repository) Health() error {
	return db.sqlDB.Ping()
}

func (db *postgres${TITLE_DOMAIN}Repository) GenerateID() (uint64, error) {
	return uint64(node.Generate().Int64()), nil
}

func (db *postgres${TITLE_DOMAIN}Repository) GenerateUUID() (string, error) {
	return uuid.NewString(), nil
}"

write_file "${DOMAIN_DIR}/delivery/http/routes-http.go" "package http

import (
	\"github.com/gofiber/fiber/v2\"

	\"${MODULE_PATH}/domain/${DOMAIN}\"
	\"${MODULE_PATH}/utils/middleware\"
)

type Http${TITLE_DOMAIN}Handler struct {
	${DOMAIN}Usecase ${DOMAIN}.Usecase
}

func NewHttp${TITLE_DOMAIN}Handler(${DOMAIN}UC ${DOMAIN}.Usecase, httpRouter *fiber.App) {
	handler := &Http${TITLE_DOMAIN}Handler{
		${DOMAIN}Usecase: ${DOMAIN}UC,
	}

	public := httpRouter.Group(\"/api/v1/${DOMAIN}\")
	public.Use(middleware.PublicMiddleware)
	public.Get(\"/health\", handler.Health)
}"

write_file "${DOMAIN_DIR}/delivery/http/http-health.go" "package http

import (
	\"${MODULE_PATH}/utils/message\"

	\"github.com/gofiber/fiber/v2\"
)

func (handler *Http${TITLE_DOMAIN}Handler) Health(ctx *fiber.Ctx) error {
	status, err := handler.${DOMAIN}Usecase.Health()
	if err != nil {
		message.ReturnInternalServerError(ctx, err.Error())
		return nil
	}

	message.ReturnOkNoParam(ctx, status)
	return nil
}"

MAIN_FILE="${BOILER_DIR}/main.go"
if [ ! -f "${MAIN_FILE}" ]; then
  echo "main.go not found in boilerplate" >&2
  exit 1
fi

IMPORT_HTTP="_${DOMAIN}HttpDeliver \"${MODULE_PATH}/domain/${DOMAIN}/delivery/http\""
IMPORT_REPO="_${DOMAIN}Repository \"${MODULE_PATH}/domain/${DOMAIN}/repository\""
IMPORT_UC="_${DOMAIN}Usecase \"${MODULE_PATH}/domain/${DOMAIN}/usecase\""

if ! grep -q "domain/${DOMAIN}/delivery/http" "${MAIN_FILE}"; then
  tmp_file="${MAIN_FILE}.tmp"
  awk -v i1="${IMPORT_HTTP}" -v i2="${IMPORT_REPO}" -v i3="${IMPORT_UC}" '
    BEGIN { inserted=0 }
    /^import \(/ && inserted==0 {
      print
      print "\t" i1
      print "\t" i2
      print "\t" i3
      inserted=1
      next
    }
    { print }
  ' "${MAIN_FILE}" > "${tmp_file}"
  mv "${tmp_file}" "${MAIN_FILE}"
fi

INIT_LINE1="${DOMAIN}Repo := _${DOMAIN}Repository.NewPostgresRepository(postgreDB)"
INIT_LINE2="${DOMAIN}UC := _${DOMAIN}Usecase.New${TITLE_DOMAIN}Usecase(${DOMAIN}Repo)"
INIT_LINE3="_${DOMAIN}HttpDeliver.NewHttp${TITLE_DOMAIN}Handler(${DOMAIN}UC, app)"

if ! grep -q "NewHttp${TITLE_DOMAIN}Handler" "${MAIN_FILE}"; then
  tmp_file="${MAIN_FILE}.tmp"
  awk -v l1="${INIT_LINE1}" -v l2="${INIT_LINE2}" -v l3="${INIT_LINE3}" '
    BEGIN { inserted=0 }
    {
      print
      if ($0 ~ /postgreDB := ConnectPostgre\(\)/ && inserted==0) {
        print ""
        print "\t" l1
        print "\t" l2
        print "\t" l3
        inserted=1
      }
    }
  ' "${MAIN_FILE}" > "${tmp_file}"
  mv "${tmp_file}" "${MAIN_FILE}"
fi

echo "Domain added to boilerplate: ${DOMAIN}"
