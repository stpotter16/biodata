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

type EntryAPI struct {
	Date   string `json:"date"`
	Weight string `json:"weight"`
	Waist  string `json:"waist"`
	BP     string `json:"bp"`
}

type Entry struct {
	Date   time.Time
	Weight Weight
	Waist  Waist
	BP     BP
}

func ToEntryApi(entry Entry) (EntryAPI, error) {
	dateStr := entry.Date.Format(time.RFC3339)
	apiEntry := EntryAPI{
		Date:   dateStr,
		Weight: entry.Weight.String(),
		Waist:  entry.Waist.String(),
		BP:     entry.BP.String(),
	}
	return apiEntry, nil
}
