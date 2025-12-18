#!/usr/bin/env bash

# Fail on first error
set -e

# Fail on unset variable
set -u

# Echo commands
set -x

SCRIPT_DIR="$( cd "$(dirname ${BASH_SOURCE[0]})" &> /dev/null && pwd)"
readonly SCRIPT_DIR
cd "${SCRIPT_DIR}/.."

go test ./...

go vet ./...

STATICCHECK_PATH="$(go env GOPATH)/bin/staticcheck"
readonly STATICCHECK_PATH
readonly STATICCHECK_VERSION="v0.6.1"
if [[ ! -f "${STATICCHECK_PATH}" ]]; then
    go install \
        -ldflags=-linkmode=external \
        "honnef.co/go/tools/cmd/staticcheck@${STATICCHECK_VERSION}"
fi

${STATICCHECK_PATH} ./...

ERRCHECK_PATH="$(go env GOPATH)/bin/errcheck"
readonly ERRCHECK_PATH
readonly ERRCHECK_VERSION="v1.9.0"
if [[ ! -f "${ERRCHECK_PATH}" ]]; then
    go install \
        -ldflags=-linkmode=external \
        "github.com/kisielk/errcheck@${ERRCHECK_VERSION}"
fi

${ERRCHECK_PATH} -ignoretests ./...
