#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./setup-project.sh --project <project_name>
  ./setup-project.sh <project_name>

Example:
  ./setup-project.sh --project svc-itasset-golang
  ./setup-project.sh backend_itasset
EOF
}

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_MOD_FILE="${ROOT_DIR}/go.mod"
README_FILE="${ROOT_DIR}/README.md"

NEW_MODULE=""
NEW_PROJECT=""
POSITIONAL_PROJECT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project)
      NEW_PROJECT="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    -*)
      echo "unknown argument: $1" >&2
      usage
      exit 1
      ;;
    *)
      if [[ -n "${NEW_PROJECT}" ]]; then
        echo "project name already provided via --project" >&2
        usage
        exit 1
      fi
      if [[ -n "${POSITIONAL_PROJECT}" ]]; then
        echo "only one project name is allowed" >&2
        usage
        exit 1
      fi
      POSITIONAL_PROJECT="$1"
      shift
      ;;
  esac
done

if [[ -z "${NEW_PROJECT}" && -n "${POSITIONAL_PROJECT}" ]]; then
  NEW_PROJECT="${POSITIONAL_PROJECT}"
fi

if [[ -z "${NEW_PROJECT}" ]]; then
  echo "project name is required" >&2
  usage
  exit 1
fi

if [[ ! "${NEW_PROJECT}" =~ ^[a-z0-9][a-z0-9_-]*$ ]]; then
  echo "invalid project name: use lowercase letters, numbers, '-' or '_'" >&2
  exit 1
fi

NEW_MODULE="${NEW_PROJECT}"

if [[ ! -f "${GO_MOD_FILE}" ]]; then
  echo "go.mod not found in ${ROOT_DIR}" >&2
  exit 1
fi

OLD_MODULE="$(awk '$1=="module"{print $2; exit}' "${GO_MOD_FILE}")"
if [[ -z "${OLD_MODULE}" ]]; then
  echo "failed to read current module from go.mod" >&2
  exit 1
fi

tmp_go_mod="${GO_MOD_FILE}.tmp"
awk -v m="${NEW_MODULE}" '
  BEGIN { changed=0 }
  $1=="module" && changed==0 {
    print "module " m
    changed=1
    next
  }
  { print }
  END {
    if (changed==0) {
      exit 1
    }
  }
' "${GO_MOD_FILE}" > "${tmp_go_mod}"
mv "${tmp_go_mod}" "${GO_MOD_FILE}"

if [[ "${OLD_MODULE}" != "${NEW_MODULE}" ]]; then
  while IFS= read -r -d '' file; do
    OLD_MODULE="${OLD_MODULE}" NEW_MODULE="${NEW_MODULE}" perl -0pi -e 's/\Q$ENV{OLD_MODULE}\E/$ENV{NEW_MODULE}/g' "${file}"
  done < <(
    find "${ROOT_DIR}" -type f \
      \( -name "*.go" -o -name "*.md" -o -name "*.sh" -o -name "go.mod" \) \
      -not -path "*/.git/*" \
      -print0
  )
fi

# Normalize any stale local import prefixes in Go files, e.g. copied boilerplate imports.
prefixes_tmp="$(mktemp)"
while IFS= read -r -d '' file; do
  perl -nE '
    while (/"([^"]+)\/(?:domain|utils|entity|valueobject)(?:\/[^"]*)?"/g) {
      say $1;
    }
  ' "${file}" >> "${prefixes_tmp}"
done < <(
  find "${ROOT_DIR}" -type f -name "*.go" -not -path "*/.git/*" -print0
)

if [[ -s "${prefixes_tmp}" ]]; then
  sort -u "${prefixes_tmp}" | while IFS= read -r prefix; do
    if [[ -n "${prefix}" && "${prefix}" != "${NEW_MODULE}" ]]; then
      while IFS= read -r -d '' gofile; do
        OLD_PREFIX="${prefix}" NEW_MODULE="${NEW_MODULE}" perl -0pi -e 's/\Q$ENV{OLD_PREFIX}\E/$ENV{NEW_MODULE}/g' "${gofile}"
      done < <(
        find "${ROOT_DIR}" -type f -name "*.go" -not -path "*/.git/*" -print0
      )
    fi
  done
fi
rm -f "${prefixes_tmp}"

if [[ -f "${README_FILE}" ]]; then
  tmp_readme="${README_FILE}.tmp"
  awk -v p="${NEW_PROJECT}" '
    NR==1 && /^# / {
      print "# " p
      next
    }
    { print }
  ' "${README_FILE}" > "${tmp_readme}"
  mv "${tmp_readme}" "${README_FILE}"
fi

if command -v go >/dev/null 2>&1; then
  (cd "${ROOT_DIR}" && go mod tidy)
else
  echo "warning: 'go' command not found, skipped dependency installation" >&2
fi

echo "Project setup updated:"
echo "  module : ${OLD_MODULE} -> ${NEW_MODULE}"
echo "  project: ${NEW_PROJECT}"
echo "  deps   : go mod tidy"
