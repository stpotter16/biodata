package sessions

import (
	"bytes"
	"encoding/json"
	"log"
)

type serializedSession struct {
	Id string `json:"id"`
}

func deserializeSession(b []byte) (Session, error) {
	var ss serializedSession
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&ss); err != nil {
		log.Printf("Failed to deserialize session from json: %v", err)
		return Session{}, err
	}

	session := Session{
		ID: ss.Id,
	}

	return session, nil
}
