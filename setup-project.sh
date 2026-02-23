#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./setup-project.sh
  ./setup-project.sh --project <project_name>
  ./setup-project.sh <project_name>

Example:
  ./setup-project.sh
  ./setup-project.sh --project svc-itasset-golang
  ./setup-project.sh backend_itasset
EOF
}

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_MOD_FILE="${ROOT_DIR}/go.mod"
README_FILE="${ROOT_DIR}/README.md"

NEW_MODULE=""
NEW_PROJECT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project)
      if [[ $# -lt 2 || -z "${2:-}" ]]; then
        echo "missing value for --project" >&2
        usage
        exit 1
      fi
      if [[ -n "${NEW_PROJECT}" ]]; then
        echo "project name already provided" >&2
        usage
        exit 1
      fi
      NEW_PROJECT="${2}"
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
        echo "only one project name is allowed" >&2
        usage
        exit 1
      fi
      NEW_PROJECT="$1"
      shift
      ;;
  esac
done

if [[ -z "${NEW_PROJECT}" ]]; then
  read -r -p "Nama project: " NEW_PROJECT
fi

if [[ -z "${NEW_PROJECT}" ]]; then
  echo "project name is required" >&2
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

normalize_local_import_prefixes() {
  while IFS= read -r -d '' gofile; do
    NEW_MODULE="${NEW_MODULE}" perl -0pi -e '
      s/"[^"\n]+\/(domain|utils|entity|valueobject)(\/[^"\n]+)"/"$ENV{NEW_MODULE}\/$1$2"/g;
      s/"[^"\n]+\/(domain|utils|entity|valueobject)"/"$ENV{NEW_MODULE}\/$1"/g;
    ' "${gofile}"
  done < <(
    find "${ROOT_DIR}" -type f -name "*.go" -not -path "*/.git/*" -print0
  )
}

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
      \( -name "*.go" -o -name "*.mod" -o -name "*.sum" -o -name "*.md" -o -name "*.sh" -o -name "*.env*" \) \
      -not -path "*/.git/*" \
      -print0
  )
fi

normalize_local_import_prefixes

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

deps_status="skipped"
if command -v go >/dev/null 2>&1; then
  if (cd "${ROOT_DIR}" && go mod tidy); then
    deps_status="go mod tidy"
  else
    deps_status="go mod tidy failed (run manually)"
    echo "warning: go mod tidy failed. run 'go mod tidy' manually when network is available." >&2
  fi
else
  deps_status="go not found (skipped)"
  echo "warning: 'go' command not found, skipped dependency installation" >&2
fi

echo "Project setup updated:"
echo "  module : ${OLD_MODULE} -> ${NEW_MODULE}"
echo "  project: ${NEW_PROJECT}"
echo "  deps   : ${deps_status}"
