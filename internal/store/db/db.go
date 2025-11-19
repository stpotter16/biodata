package db

import (
	"database/sql"
	"log"

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

	writeDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Could not open write db: %v", err)
		return DB{}, err
	}

	db := DB{
		readDB:  readDB,
		writeDB: writeDB,
	}

	return db, nil
}
