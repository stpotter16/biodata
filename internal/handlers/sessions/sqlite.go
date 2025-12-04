package sessions

import (
	"context"
	"log"
	"time"
)

const CLEANUP_INTERVAL = 5 * time.Minute
const SESSION_TTL = 60 * time.Minute

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
		session
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

func (s SessionManger) insertSession(key string, session []byte) error {
	insert := `
	INSERT OR REPLACE INTO
		session
	(
		key,
		value,
		expires_at
	)
	VALUES (
		?,
		?,
		?
	)`
	expires_time := time.Now().Add(SESSION_TTL).Format(time.RFC3339)

	// TODO - what context?
	_, err := s.db.Exec(
		context.TODO(),
		insert,
		key,
		session,
		expires_time,
	)

	return err
}
