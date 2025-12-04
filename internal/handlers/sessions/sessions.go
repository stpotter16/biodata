package sessions

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/stpotter16/biodata/internal/store/db"
)

const SESSION_KEY = "session"

type Session struct {
	ID string
}

type SessionManger struct {
	db                   *db.DB
	sessionHmacSecretKey string
}

func New(db db.DB) SessionManger {
	s := SessionManger{
		db:                   &db,
		sessionHmacSecretKey: "secret", // TODO - fixme
	}

	go s.sessionCleanup()

	return s
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
