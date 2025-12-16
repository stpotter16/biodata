DROP TABLE session;

CREATE TABLE IF NOT EXISTS session (
    session_key TEXT PRIMARY KEY,
    value BLOB,
    expires_at TEXT NOT NULL
) STRICT;
