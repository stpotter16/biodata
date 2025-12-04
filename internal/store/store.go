package store

import (
	"time"

	"github.com/stpotter16/biodata/internal/types"
)

type Store interface {
	GetEntries() ([]types.Entry, error)
	GetEntry(entryDate time.Time) (types.Entry, error)
	InsertEntry(types.Entry) error
	UpdateEntry(types.Entry) error
}
