package sqlite

import (
	"log"

	"github.com/stpotter16/biodata/internal/store/db"
)

type Store struct {
	db *db.DB
}

func New(db db.DB) (Store, error) {
	store := Store{db: &db}
	err := store.runMigrations()
	if err != nil {
		log.Printf("Could not run database migrations: %v", err)
		return Store{}, err
	}
	return store, nil
}
