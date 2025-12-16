CREATE TABLE IF NOT EXISTS entry2 (
    id INTEGER PRIMARY KEY,
    entry_date TEXT NOT NULL UNIQUE,
    weight REAL,
    waist REAL,
    bp TEXT,
    created TEXT NOT NULL,
    last_modified TEXT NOT NULL
) STRICT;

INSERT INTO entry2
SELECT
    id,
    date AS entry_date,
    weight,
    waist,
    bp,
    created,
    last_modified
FROM
    entry;

DROP TABLE entry;

ALTER TABLE entry2
RENAME TO entry;

CREATE INDEX IF NOT EXISTS idx_entry_date ON entry (entry_date);
