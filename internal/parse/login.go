package parse

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stpotter16/biodata/internal/types"
)

func ParseLoginPost(r *http.Request) (types.LoginRequest, error) {
	body := struct {
		Passphrase string `json:"passphrase"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		log.Printf("Invalid login request: %v", err)
		return types.LoginRequest{}, nil
	}

	request := types.LoginRequest{
		Passphrase: body.Passphrase,
	}
	return request, nil
}
