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
const SESSION_ENV_KEY = "BIODATA_SESSION_ENV_KEY"

const DEFAULT_USER_ID = 1

type Session struct {
	ID     string
	UserId uint8
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
		ID:     sessionId,
		UserId: DEFAULT_USER_ID,
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

func (s SessionManger) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.loadSession(r)
	if err != nil {
		log.Printf("Could not load session: %v", err)
		return err
	}
	if err = s.deleteSession(session.ID); err != nil {
		log.Printf("Could not delete session: %v", err)
		return err
	}
	if err = s.deleteSessionCookie(w); err != nil {
		log.Printf("Could not delete session cookie: %v", err)
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

func (s SessionManger) SessionFromContext(ctx context.Context) (Session, error) {
	session, okay := ctx.Value(SESSION_KEY).(Session)
	if !okay {
		log.Printf("Unable to extract session from context")
		return Session{}, errors.New("No session info in context")
	}
	return session, nil
}

func (s SessionManger) loadSession(r *http.Request) (Session, error) {
	cookie, err := s.readSessionCookie(r)
	if err != nil {
		log.Printf("Failed to read session cookie: %v", err)
		return Session{}, err
	}

	cookieVals := strings.SplitN(cookie, "::", 2)
	if len(cookieVals) != 2 {
		log.Printf("Invalid cookie value: %s", cookie)
		return Session{}, errors.New("Cookie is invalid")
	}
	cookieToken := cookieVals[1]

	serializedSession, err := s.readSession(cookieToken)
	if err != nil {
		log.Printf("Failed to load session data: %v", err)
		return Session{}, err
	}

	session, err := deserializeSession(serializedSession)
	if err != nil {
		return Session{}, err
	}

	if cookieToken != session.ID {
		return Session{}, errors.New("Invalid session token")
	}

	return session, nil
}
