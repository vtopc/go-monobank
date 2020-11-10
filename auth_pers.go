package monobank

import (
	"net/http"
)

const (
	personalTokenKey = "X-Token"
)

type PersAuth struct {
	token string
}

func (a PersAuth) SetAuth(req *http.Request) error {
	if req == nil {
		return nil
	}

	req.Header.Set(personalTokenKey, a.token)

	return nil
}

func NewPersonalAuthorizer(token string) PersAuth {
	return PersAuth{token: token}
}
