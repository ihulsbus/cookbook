#!/usr/bin/env bash
# Copies the repo's single source-of-truth OpenAPI spec (../openapi.yaml) into
# the docs Helm chart's files/ directory so it can be packaged into the
# ConfigMap (see k8s/templates/configmap.yaml). Helm cannot reference files
# outside the chart directory, so this copy must be refreshed before
# `helm package` / `helm upgrade --install` whenever openapi.yaml changes.
#
# Usage: ./docs/sync-openapi.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

mkdir -p "${SCRIPT_DIR}/k8s/files"
cp "${REPO_ROOT}/openapi.yaml" "${SCRIPT_DIR}/k8s/files/openapi.yaml"

echo "Synced ${REPO_ROOT}/openapi.yaml -> ${SCRIPT_DIR}/k8s/files/openapi.yaml"
