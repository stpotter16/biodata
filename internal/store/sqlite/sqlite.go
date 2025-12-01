package sqlite

import (
	"log"
	"time"

	"github.com/stpotter16/biodata/internal/store/db"
)

const (
	timeFormat = time.RFC3339
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

func formatTime(t time.Time) string {
	return t.Format(timeFormat)
}
