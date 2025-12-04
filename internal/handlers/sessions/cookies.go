package sessions

import (
	"net/http"

	"github.com/stpotter16/biodata/internal/cookies"
)

const SESSION_COOKIE = "X-BIODATA-SESSION"

func (s SessionManger) readSessionCookie(r *http.Request) (string, error) {
	return cookies.ReadSigned(r, SESSION_COOKIE, s.sessionHmacSecretKey)
}
