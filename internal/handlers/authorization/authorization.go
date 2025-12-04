package authorization

import "errors"

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

func (a Authorizer) Authorize(passphrase string) bool {
	if passphrase != a.passphrase {
		return false
	}
	return true
}
