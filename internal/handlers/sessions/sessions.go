package sessions

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/stpotter16/biodata/internal/store/db"
)

const SESSION_KEY = "session"
const SESSION_ENV_KEY = "SESSION_ENV_KEY"

type Session struct {
	ID string
}

type SessionManger struct {
	db                   *db.DB
	sessionHmacSecretKey string
}

func New(db db.DB, getenv func(string) string) (SessionManger, error) {
	hmacSecret := getenv(SESSION_ENV_KEY)
	if hmacSecret == "" {
		return SessionManger{}, errors.New("Could not locate HMAC secret key")
	}

	s := SessionManger{
		db:                   &db,
		sessionHmacSecretKey: hmacSecret,
	}

	go s.sessionCleanup()

	return s, nil
}

func (s SessionManger) CreateSession(w http.ResponseWriter) error {
	sessionId := uuid.NewString()
	session := Session{
		ID: sessionId,
	}

	if err := s.writeSessionCookie(w, session); err != nil {
		log.Printf("Failed to set session cookie: %v", err)
		return err
	}

	serializedSession, err := serializeSession(session)
	if err != nil {
		log.Printf("Failed to serialize session: %v", err)
		return err
	}

	if err := s.insertSession(session.ID, serializedSession); err != nil {
		log.Printf("Failed to save session: %v", err)
		return err
	}
	return nil
}

func (s SessionManger) PopulateSessionContext(r *http.Request) (context.Context, error) {
	session, err := s.loadSession(r)

	if err != nil {
		log.Printf("Unable to populate session context: %v", err)
		return nil, err
	}

	return context.WithValue(r.Context(), SESSION_KEY, session), nil
}

func (s SessionManger) loadSession(r *http.Request) (Session, error) {
	cookie, err := s.readSessionCookie(r)
	if err != nil {
		log.Printf("Failed to read session cookie: %v", err)
		return Session{}, err
	}

	cookieVals := strings.SplitN(cookie, "::", 2)
	cookieKey := cookieVals[0]
	cookieToken := cookieVals[1]

	serializedSession, err := s.readSession(cookieKey)
	if err != nil {
		log.Printf("Failed to load session data: %v", err)
		return Session{}, nil
	}

	session, err := deserializeSession(serializedSession)
	if err != nil {
		return Session{}, nil
	}

	if cookieToken != session.ID {
		return Session{}, errors.New("Invalid session token")
	}

	return session, nil
}
