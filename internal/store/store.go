package store

import (
	"context"
	"time"

	"github.com/stpotter16/biodata/internal/types"
)

type Store interface {
	GetEntries(ctx context.Context) ([]types.Entry, error)
	GetLastTenEntries(ctx context.Context) ([]types.Entry, error)
	GetEntry(ctx context.Context, entryDate time.Time) (types.Entry, error)
	InsertEntry(ctx context.Context, entry types.Entry) error
	UpdateEntry(ctx context.Context, entry types.Entry) error
}
