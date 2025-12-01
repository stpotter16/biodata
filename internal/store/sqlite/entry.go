package sqlite

import (
	"context"
	"time"

	"github.com/stpotter16/biodata/internal/parse"
	"github.com/stpotter16/biodata/internal/types"
)

func (s Store) GetEntries() ([]types.Entry, error) {
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

	var entries []types.Entry
	for rows.Next() {
		var e types.EntryDTO
		var date string
		err := rows.Scan(&e.Id, &date, &e.Weight, &e.Waist, &e.Bp, &e.Created, &e.LastModified)
		if err != nil {
			return nil, err
		}
		entryDate, err := parseTime(date)
		if err != nil {
			return nil, err
		}
		e.Date = entryDate
		entry, err := parse.ParseEntryDTO(e)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
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
	formatedDate := formatTime(entry.Date)

	// TODO - what context?
	_, err := s.db.Exec(
		context.TODO(),
		insert,
		formatedDate,
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
