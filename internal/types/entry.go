package types

import "database/sql"

type EntryDTO struct {
	Id           int
	Date         string
	Weight       sql.NullFloat64
	Waist        sql.NullFloat64
	Bp           sql.NullString
	Created      string
	LastModified string
}
