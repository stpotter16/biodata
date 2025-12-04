package sessions

import (
	"context"
	"log"
	"time"
)

const CLEANUP_INTERVAL = 5 * time.Minute

func (s SessionManger) sessionCleanup() {
	ticker := time.NewTicker(CLEANUP_INTERVAL)

	for range ticker.C {
		if err := s.deleteExpiredSessions(); err != nil {
			log.Printf("Failed to delete expired sessions: %v", err)
		}
	}
}

func (s SessionManger) deleteExpiredSessions() error {
	delete := `
	DELETE FROM
		sessions
	WHERE
		expires_at <= datetime('now', 'localtime')
	`
	// TODO - what context?
	_, err := s.db.Exec(context.TODO(), delete)

	return err
}
