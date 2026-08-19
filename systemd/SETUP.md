# systemd setup: biodata + litestream

Runs the `biodata` binary as a dedicated unprivileged user, with litestream
(S3 replication) also running unprivileged as its own user, sharing access
to the data directory via a common group.

## 1. Build the binary

```bash
make server/release
```

## 2. Create users and the shared group

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin biodata
sudo useradd --system --no-create-home --shell /usr/sbin/nologin litestream

sudo groupadd biodata-data
sudo usermod -aG biodata-data biodata
sudo usermod -aG biodata-data litestream
```

## 3. Create directories and set permissions

```bash
sudo mkdir -p /opt/biodata /etc/biodata /var/lib/biodata

sudo chown biodata:biodata-data /var/lib/biodata
sudo chmod 2770 /var/lib/biodata   # setgid: new files inherit the biodata-data group
```

## 4. Install the biodata binary and unit

```bash
sudo cp release/biodata /opt/biodata/biodata
sudo cp systemd/biodata.service /etc/systemd/system/biodata.service

sudo cp systemd/biodata.env.example /etc/biodata/biodata.env
sudo chown root:biodata /etc/biodata/biodata.env
sudo chmod 640 /etc/biodata/biodata.env
sudo $EDITOR /etc/biodata/biodata.env   # fill BIODATA_PASSPHRASE, BIODATA_SESSION_ENV_KEY
```

`/etc/biodata/biodata.env` should contain:

```
PORT=8080
BIODATA_DB_PATH=/var/lib/biodata
BIODATA_PASSPHRASE=<generated>
BIODATA_SESSION_ENV_KEY=<generated>
```

## 5. Override litestream's packaged unit to drop root

litestream's out-of-box systemd unit has no `User=`, so it runs as root by
default, and has no `[Unit]` ordering at all — no `After=`/
`Wants=network-online.target`, unlike `biodata.service`. Override both:

```bash
sudo mkdir -p /etc/systemd/system/litestream.service.d
sudo tee /etc/systemd/system/litestream.service.d/override.conf > /dev/null <<'EOF'
[Unit]
After=network-online.target
Wants=network-online.target

[Service]
User=litestream
Group=litestream
EOF
```

## 6. Install litestream.yml from the repo

Always deploy `litestream.yml` from this repo, not a hand-edited copy on
the server — a stale copy left over from before the `biotrak` → `biodata`
rename is exactly what caused the "bucket required for s3 replica" error
during initial setup.

```bash
sudo cp litestream.yml /etc/litestream.yml
sudo chown root:litestream /etc/litestream.yml
sudo chmod 640 /etc/litestream.yml
```

## 7. Secure litestream's own secrets (S3 credentials + db path)

`litestream.yml` pulls its values from env vars — `${ACCESS_KEY_ID}`,
`${SECRET_ACCESS_KEY}`, `${BIOTRAK_BUCKET}`, `${BIOTRAK_ENDPOINT}`, and
`${BIOTRAK_PATH}` (the db path — note this is a *different* variable name
than the biodata app's own `BIODATA_DB_PATH` in `/etc/biodata/biodata.env`,
since these are two separate services/env files; set it to the same
directory value, `/var/lib/biodata`). If the packaged install doesn't
already give you a hook for supplying these, add one:

```bash
sudo mkdir -p /etc/litestream
sudo tee /etc/litestream/litestream.env > /dev/null <<'EOF'
ACCESS_KEY_ID=
SECRET_ACCESS_KEY=
BIOTRAK_BUCKET=
BIOTRAK_ENDPOINT=
BIOTRAK_PATH=/var/lib/biodata
EOF
sudo chown root:litestream /etc/litestream/litestream.env
sudo chmod 640 /etc/litestream/litestream.env
sudo $EDITOR /etc/litestream/litestream.env   # fill in real values

# add EnvironmentFile to the override from step 5 (rewrite the whole file —
# drop-ins are per-file, not per-line, so include everything each time)
sudo tee /etc/systemd/system/litestream.service.d/override.conf > /dev/null <<'EOF'
[Unit]
After=network-online.target
Wants=network-online.target

[Service]
User=litestream
Group=litestream
EnvironmentFile=/etc/litestream/litestream.env
EOF
```

## 8. Restore data (manual, first boot / disaster recovery only)

`litestream.service` only runs `litestream replicate` — there's no
automatic restore-on-start. That's deliberate: restore only ever matters
once per box (initial provisioning, or actual disaster recovery), so it
doesn't belong in the normal, unattended start path. Baking it into
`ExecStartPre=` made every ordinary `systemctl restart litestream` (or a
reboot) pay the cost of `litestream restore`'s generation/WAL enumeration —
which can take minutes once a replica has accumulated many generations —
and, because `biodata.service` was set to require litestream, a slow or
failed restore could take biodata down with it too. Run this by hand
instead, once, before starting either service for the first time on a
box (skip it entirely for a genuinely brand-new deployment with no
existing backup — biodata will just create a fresh database):

```bash
sudo -u litestream bash -c '
  set -a; source /etc/litestream/litestream.env; set +a
  litestream restore -v -if-replica-exists "${BIOTRAK_PATH}/biodata.sqlite"
'
ls -la /var/lib/biodata/
```

litestream writes restored files as `-rw-------` (owner `litestream` only)
regardless of umask — this is not the same gap `UMask=0007` on
`biodata.service` fixes (step 3/Notes), since that only covers files the
`biodata` process itself creates, not litestream's. Left alone, biodata
can't open the db it just restored. Fix the mode before starting
`biodata.service`, and clean up litestream's leftover restore scratch
files while you're at it (harmless, but no reason to keep them —
`biodata` creates its own `-wal`/`-shm` at the real path when it opens
the db):

```bash
sudo chmod 660 /var/lib/biodata/biodata.sqlite
sudo rm -f /var/lib/biodata/biodata.sqlite.tmp-wal /var/lib/biodata/biodata.sqlite.tmp-shm
ls -la /var/lib/biodata/
```

`systemd/litestream-restore-if-missing.sh` in the repo does the same
"only restore if the local file doesn't already exist" check, plus the
`chmod`/cleanup above automatically, as a convenience wrapper — copy it to
the box and run it directly if you'd rather not do those steps by hand:

```bash
sudo cp systemd/litestream-restore-if-missing.sh /usr/local/bin/litestream-restore-if-missing.sh
sudo chown root:root /usr/local/bin/litestream-restore-if-missing.sh
sudo chmod 755 /usr/local/bin/litestream-restore-if-missing.sh

sudo -u litestream bash -c '
  set -a; source /etc/litestream/litestream.env; set +a
  /usr/local/bin/litestream-restore-if-missing.sh
'
```

If a restore ever does seem to hang, time it directly rather than
guessing — note `restore` looks up its replica config by matching its
`DB_PATH` argument against `litestream.yml`'s configured `dbs[].path`, so
pass the real path and use `-o` to redirect the actual output elsewhere
to test safely without touching the real db file:

```bash
sudo -u litestream bash -c '
  set -a; source /etc/litestream/litestream.env; set +a
  time litestream restore -o /tmp/biodata-restore-test.sqlite -if-replica-exists "${BIOTRAK_PATH}/biodata.sqlite"
'
ls -la /tmp/biodata-restore-test.sqlite
sudo rm -f /tmp/biodata-restore-test.sqlite
```

## 9. Start everything

```bash
sudo systemctl daemon-reload

sudo systemctl enable --now biodata.service
sudo systemctl restart litestream
sudo systemctl enable litestream   # if not already enabled by the package

sudo systemctl status biodata.service litestream
sudo systemctl cat litestream       # confirm the override merged (User=litestream, Group=litestream)
journalctl -u biodata.service -u litestream -f
```

## 10. Verify

```bash
sudo -u litestream test -r /var/lib/biodata/biodata.sqlite && echo "litestream can read the db"
sudo -u litestream test -w /var/lib/biodata/biodata.sqlite && echo "litestream can write to the db"
sudo -u biodata test -w /var/lib/biodata && echo "biodata can write to the data dir"
```

If the second check fails (or `journalctl -u litestream` shows "attempt to
write to a readonly database"), see the UMask note below — it's almost
always the actual db/WAL files being group-unwritable, not the directory.

## Fallback to root if litestream can't start

```bash
sudo rm -rf /etc/systemd/system/litestream.service.d
sudo systemctl daemon-reload
sudo systemctl restart litestream
```

Root already has full access to `/var/lib/biodata`, so no group changes need
reverting — `biodata.service`'s own hardening is a separate unit and
unaffected either way.

## Notes

- **Group-writable db files (UMask)**: the setgid bit on `/var/lib/biodata`
  (`chmod 2770` in step 3) makes new files inherit the `biodata-data`
  *group*, but it doesn't change the *mode* bits — with systemd's default
  `umask 0022`, files biodata creates come out `rw-r--r--`, so litestream
  (a group member, not the owner) can read `biodata.sqlite` but not write
  to it, causing "attempt to write to a readonly database" in
  `journalctl -u litestream`. `biodata.service` sets `UMask=0007` so any
  file biodata creates (the db, and the `-wal`/`-shm` files SQLite creates
  alongside it in WAL mode) comes out `rw-rw----` instead. If you hit this
  on an existing install, fix the already-created files directly:
  `sudo chmod 660 /var/lib/biodata/biodata.sqlite*`. Don't chmod the
  *directory* to `660` — directories need the execute bit to be
  traversable at all, so that would break both processes' access
  entirely; `2770` from step 3 is already correct.
- **litestream.yml vars vs .env.example**: `litestream.yml`'s env var names
  (`ACCESS_KEY_ID`, `SECRET_ACCESS_KEY`, `BIOTRAK_BUCKET`,
  `BIOTRAK_ENDPOINT`, `BIOTRAK_PATH`) don't match `.env.example`'s
  `LITESTREAM_*` names. This divergence predates the systemd deployment
  and hasn't been reconciled yet — pick whichever naming wins and update
  both when touching this.
- **SIGINT vs SIGTERM**: `cmd/server/main.go` only handles `os.Interrupt`
  (SIGINT), not SIGTERM. `biodata.service` sets `KillSignal=SIGINT` so
  `systemctl stop`/`restart` still trigger the app's graceful
  `server.Shutdown()` instead of a hard kill.
- **Startup ordering**: `biodata.service` sets `After=litestream.service`
  and `Wants=litestream.service` (a soft dependency — `systemctl start
  biodata.service` also brings up litestream, but a litestream failure
  won't take biodata down with it). This used to be `Requires=`, back when
  litestream's `ExecStartPre=` ran restore automatically and biodata
  genuinely needed that to finish first; now that restore is a manual
  step (see step 8), there's no automatic ordering dependency to enforce
  and the hard requirement was just extra blast radius for no benefit.
- **`snapshot-interval` and restore time**: litestream 0.3.x (what's
  pinned here) has no compaction — unlike the
  newer 0.5.x LTX rewrite, it never merges old WAL segment objects into
  larger ones. WAL segments just accumulate between snapshots, so a long
  `snapshot-interval` combined with a short `sync-interval` (`1h` and
  `1s` respectively here) can build up a large number of small objects
  that a future restore has to enumerate — directly increasing restore
  time as a generation ages. `litestream.yml` sets `snapshot-interval: 1h`
  (litestream's own docs suggest this trade-off for exactly this reason)
  to bound the worst case, at the cost of more frequent full-db snapshot
  storage. If restores start feeling slow again as real usage grows,
  check how many objects have accumulated in the bucket before assuming
  something's broken — it may just be that this interval needs to come
  down further, or old generations need cleaning up.
