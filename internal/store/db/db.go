package db

import (
	"context"
	"database/sql"
	"log"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	writeDB *sql.DB
	readDB  *sql.DB
}

func New(dbPath string) (DB, error) {
	readDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Could not open read db: %v", err)
		return DB{}, err
	}
	readDB.SetMaxOpenConns(max(4, runtime.NumCPU()))
	if err = applyPragmas(readDB); err != nil {
		log.Printf("Could not apply PRAGMAs to read db: %v", err)
		return DB{}, err
	}

	writeDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Could not open write db: %v", err)
		return DB{}, err
	}
	writeDB.SetMaxOpenConns(1)
	if err = applyPragmas(writeDB); err != nil {
		log.Printf("Could not apply PRAGMAs to write db: %v", err)
		return DB{}, err
	}

	db := DB{
		readDB:  readDB,
		writeDB: writeDB,
	}

	return db, nil
}

func (db DB) QueryRow(ctx context.Context, dest any, query string, args ...any) error {
	row := db.readDB.QueryRowContext(ctx, query, args...)
	if err := row.Scan(dest); err != nil {
		return err
	}
	return nil
}

func (db DB) ExecuteTransaction(ctx context.Context, transactions ...string) error {
	tx, err := db.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, statement := range transactions {
		_, err = tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func applyPragmas(db *sql.DB) error {
	// TODO - maybe use context here
	if _, err := db.Exec(`
		-- https://litestream.io/tips/
		-- https://kerkour.com/sqlite-for-servers
		PRAGMA journal_mode = WAL;
		PRAGMA busy_timeout = 5000;
		PRAGMA synchronous = NORMAL;
        PRAGMA wal_autocheckpoint = 0;
		PRAGMA cache_size = 1000000000;
		PRAGMA foreign_keys = true;
		PRAGMA temp_store = memory;
	`); err != nil {
		return err
	}
	return nil
}
