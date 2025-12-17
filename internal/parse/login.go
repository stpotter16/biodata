package parse

import (
	"encoding/json"
	"net/http"

	"github.com/stpotter16/biodata/internal/types"
)

func ParseLoginPost(r *http.Request) (types.LoginRequest, error) {
	body := struct {
		Passphrase string `json:"passphrase"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		return types.LoginRequest{}, err
	}

	request := types.LoginRequest{
		Passphrase: body.Passphrase,
	}
	return request, nil
}
