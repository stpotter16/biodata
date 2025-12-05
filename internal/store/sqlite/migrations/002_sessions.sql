CREATE TABLE IF NOT EXISTS session (
    key INTEGER PRIMARY KEY,
    value BLOB,
    expires_at TEXT NOT NULL
);
