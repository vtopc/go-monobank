package monobank

import (
	"net/http"
)

const (
	personalTokenKey = "X-Token"
)

type PersonalAuthorizer struct {
	token string
}

func (a PersonalAuthorizer) SetAuth(req *http.Request) error {
	if req == nil {
		return nil
	}

	req.Header.Set(personalTokenKey, a.token)

	return nil
}

func NewPersonalAuthorizer(token string) PersonalAuthorizer {
	return PersonalAuthorizer{token: token}
}
