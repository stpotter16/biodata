package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/stpotter16/biodata/internal/parse"
	"github.com/stpotter16/biodata/internal/types"
)

func (s Store) GetEntries() ([]types.Entry, error) {
	query := `
	SELECT id, entry_date, weight, waist, bp, created, last_modified
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

func (s Store) GetEntry(entryDate time.Time) (types.Entry, error) {
	query := `
	SELECT id, entry_date, weight, waist, bp, created, last_modified
	FROM entry
	WHERE entry_date = ?
	`

	dateStr := formatTime(entryDate)
	// TODO - what to do with this context
	row := s.db.QueryRow(context.TODO(), query, dateStr)
	var entryDTO types.EntryDTO
	var date string

	if err := row.Scan(&entryDTO.Id, &date, &entryDTO.Weight, &entryDTO.Waist, &entryDTO.Bp, &entryDTO.Created, &entryDTO.LastModified); err != nil {
		return types.Entry{}, err
	}
	entryDate, err := parseTime(date)
	if err != nil {
		return types.Entry{}, err
	}
	entryDTO.Date = entryDate
	entry, err := parse.ParseEntryDTO(entryDTO)
	if err != nil {
		return types.Entry{}, err
	}

	return entry, nil
}

func (s Store) InsertEntry(entry types.Entry) error {
	insert := `
	INSERT INTO entry
	(entry_date, weight, waist, bp, created, last_modified)
	VALUES (?, ?, ?, ?, ?, ?);
	`
	now := formatTime(time.Now())
	formatedDate := formatTime(entry.Date)

	var weight sql.NullFloat64
	if entry.Weight.Valid() {
		weight = sql.NullFloat64{
			Float64: *entry.Weight.Value,
			Valid:   true,
		}
	} else {
		weight = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}

	var waist sql.NullFloat64
	if entry.Waist.Valid() {
		waist = sql.NullFloat64{
			Float64: *entry.Waist.Value,
			Valid:   true,
		}
	} else {
		waist = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}

	var bp sql.NullString
	if entry.BP.Valid() {
		bp = sql.NullString{
			String: entry.BP.String(),
			Valid:  true,
		}
	} else {
		bp = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	// TODO - what context?
	_, err := s.db.Exec(
		context.TODO(),
		insert,
		formatedDate,
		weight,
		waist,
		bp,
		now,
		now,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) UpdateEntry(entry types.Entry) error {
	update := `
	UPDATE entry
	SET weight = ?, waist = ?, bp = ?, last_modified = ?
	WHERE entry_date = ?
	`
	now := formatTime(time.Now())
	formatedTime := formatTime(entry.Date)

	var weight sql.NullFloat64
	if entry.Weight.Valid() {
		weight = sql.NullFloat64{
			Float64: *entry.Weight.Value,
			Valid:   true,
		}
	} else {
		weight = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}

	var waist sql.NullFloat64
	if entry.Waist.Valid() {
		waist = sql.NullFloat64{
			Float64: *entry.Waist.Value,
			Valid:   true,
		}
	} else {
		waist = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}

	var bp sql.NullString
	if entry.BP.Valid() {
		bp = sql.NullString{
			String: entry.BP.String(),
			Valid:  true,
		}
	} else {
		bp = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	// TODO - what context?
	_, err := s.db.Exec(
		context.TODO(),
		update,
		weight,
		waist,
		bp,
		now,
		formatedTime,
	)

	if err != nil {
		return err
	}

	return nil
}
