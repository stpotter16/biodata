package sqlite

import (
	"context"
	"time"

	"github.com/stpotter16/biodata/internal/types"
)

func (s Store) GetEntries() ([]types.EntryDTO, error) {
	query := `
	SELECT id, date, weight, waist, bp, created, last_modified
	FROM entry
	ORDER BY id DESC;
	`

	// TODO - what to do with this context
	rows, err := s.db.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []types.EntryDTO
	for rows.Next() {
		var e types.EntryDTO
		err := rows.Scan(&e.Id, &e.Date, &e.Weight, &e.Waist, &e.Bp, &e.Created, &e.LastModified)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s Store) InsertEntry(entry types.Entry) error {
	insert := `
	INSERT INTO entry
	(date, weight, waist, bp, created, last_modified)
	VALUES (?, ?, ?, ?, ?, ?);
	`
	now := formatTime(time.Now())

	// TODO - what context?
	_, err := s.db.Exec(
		context.TODO(),
		insert,
		entry.Date,
		entry.Weight,
		entry.Waist,
		entry.BP,
		now,
		now,
	)
	if err != nil {
		return err
	}
	return nil
}
