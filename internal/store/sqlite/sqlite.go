package sqlite

import "github.com/stpotter16/biodata/internal/store/db"

type Store struct {
	db *db.DB
}

func New(db db.DB) Store {
	return Store{db: &db}
}
