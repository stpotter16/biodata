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

	writeDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Could not open write db: %v", err)
		return DB{}, err
	}
	writeDB.SetMaxOpenConns(1)

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
