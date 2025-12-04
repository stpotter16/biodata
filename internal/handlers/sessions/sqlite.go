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

func (s SessionManger) readSession(key string) ([]byte, error) {
	query := `
	SELECT
		value
	FROM
		session
	WHERE
		key = ? AND
		expires_at >= datetime('now', 'localtime')
	`

	// TODO - what context
	row := s.db.QueryRow(context.TODO(), query, key)
	var serializedSession []byte
	if err := row.Scan(&serializedSession); err != nil {
		return nil, err
	}
	return serializedSession, nil
}
