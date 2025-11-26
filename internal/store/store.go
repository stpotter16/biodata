package store

import "github.com/stpotter16/biodata/internal/types"

type Store interface {
	GetEntries() ([]types.EntryDTO, error)
}
