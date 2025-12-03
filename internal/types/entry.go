package types

import (
	"database/sql"
	"time"
)

type EntryDTO struct {
	Id           int
	Date         time.Time
	Weight       sql.NullFloat64
	Waist        sql.NullFloat64
	Bp           sql.NullString
	Created      string
	LastModified string
}

type Entry struct {
	Date   time.Time
	Weight Weight
	Waist  Waist
	BP     BP
}
