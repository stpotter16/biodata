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

# Env variables for testing
SESSION_ENV_KEY="$(xxd -l32 /dev/urandom | xxd -r -ps | base64 | tr -d = | tr + - | tr / _)"
PASSPHRASE="e2epassphrase"
DB="$(mktemp -d)/biodata"

cleanup() {
    echo "Cleaning up go server PID: ${GOPID}"
    kill "${GOPID}"
}

trap cleanup EXIT INT

# Build server and run it
./dev-scripts/build-server.sh
BIODATA_SESSION_ENV_KEY=$SESSION_ENV_KEY BIODATA_PASSPHRASE=$PASSPHRASE PORT=8081 BIODATA_DB_PATH=$DB ./tmp/server &

GOPID=$!

PASSPHRASE=$PASSPHRASE npx playwright test
