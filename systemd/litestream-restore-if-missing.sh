#!/usr/bin/env bash

# Restore biodata.sqlite from the litestream replica if it doesn't already
# exist locally. Run by hand as a one-time step (initial provisioning /
# disaster recovery) per systemd/SETUP.md — not wired into any systemd unit.

set -euo pipefail

DB_FILE="${BIOTRAK_PATH}/biodata.sqlite"

if [[ -f "${DB_FILE}" ]]; then
  echo "Existing database found at ${DB_FILE} ($(stat -c %s "${DB_FILE}") bytes), skipping restore"
else
  echo "No existing database found at ${DB_FILE}, attempting restore from replica"
  litestream restore -if-replica-exists "${DB_FILE}"

  if [[ -f "${DB_FILE}" ]]; then
    # litestream writes restored files as owner-only (litestream:biodata-data,
    # mode 600) regardless of umask, which biodata.service's own UMask=0007
    # fix doesn't cover — that only applies to files biodata itself creates.
    chmod 660 "${DB_FILE}"
    rm -f "${DB_FILE}.tmp-wal" "${DB_FILE}.tmp-shm"
  fi
fi
