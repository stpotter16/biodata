package sessions

import "github.com/stpotter16/biodata/internal/store/db"

type SessionManger struct {
	db *db.DB
}

func New(db db.DB) SessionManger {
	s := SessionManger{
		db: &db,
	}

	//go s.sessionCleanup()

	return s
}
