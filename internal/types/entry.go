package types

import (
	"database/sql"
	"time"
)

type EntryDTO struct {
	Id           int
	Date         string
	Weight       sql.NullFloat64
	Waist        sql.NullFloat64
	Bp           sql.NullString
	Created      string
	LastModified string
}

type Entry struct {
	Date   time.Time
	Weight float64
	Waist  float64
	BP     string // TODO - this could be it's own type
}
