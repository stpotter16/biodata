package authorization

import (
	"errors"

	"github.com/stpotter16/biodata/internal/types"
)

type Authorizer struct {
	passphrase string
}

func New(getenv func(string) string) (Authorizer, error) {
	passphrase := getenv("BIODATA_PASSPHRASE")
	if passphrase == "" {
		return Authorizer{}, errors.New("Could not locate passphrase environment variable")
	}

	a := Authorizer{
		passphrase: passphrase,
	}

	return a, nil
}

func (a Authorizer) Authorize(loginRequest types.LoginRequest) bool {
	if loginRequest.Passphrase == a.passphrase {
		return true
	}
	return false
}

func (a Authorizer) AuthorizeApi(headerVal string) bool {
	if headerVal == a.passphrase {
		return true
	}
	return false
}
