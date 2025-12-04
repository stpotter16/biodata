CREATE TABLE IF NOT EXISTS session (
    key TEXT PRIMARY KEY,
    value BLOB,
    expires_at TEXT NOT NULL
);
