package store

import "github.com/stpotter16/biodata/internal/types"

type Store interface {
	GetEntries() ([]types.Entry, error)
	InsertEntry(types.Entry) error
}
