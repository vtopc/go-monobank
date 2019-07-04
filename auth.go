package monobank

import (
	"net/http"
)

type Authorizer interface {
	// SetAuth modifies http.Request and sets authorization tokens
	SetAuth(*http.Request)
}

const (
	personalTokenKey = "X-Token"
)

type PersonalAuthorizer struct {
	token string
}

func (a PersonalAuthorizer) SetAuth(req *http.Request) {
	if req == nil {
		return
	}

	req.Header.Set(personalTokenKey, a.token)
}

func NewPersonalAuthorizer(token string) PersonalAuthorizer {
	return PersonalAuthorizer{token}
}

// TODO: add CorporateAuthorizer https://api.monobank.ua/docs/corporate.html
