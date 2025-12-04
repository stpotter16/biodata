CREATE TABLE IF NOT EXISTS entry (
    id INTEGER PRIMARY KEY,
    date TEXT NOT NULL UNIQUE,
    weight REAL,
    waist REAL,
    bp TEXT,
    created TEXT NOT NULL,
    last_modified TEXT NOT NULL
) STRICT;

CREATE INDEX IF NOT EXISTS idx_entry_date ON entry(date);
