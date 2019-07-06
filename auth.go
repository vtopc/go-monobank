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

type NoopAuthorizer struct{}

func (a NoopAuthorizer) SetAuth(req *http.Request) {}

func NewNoopAuthorizer() NoopAuthorizer {
	return NoopAuthorizer{}
}

// TODO: add NewCorporateAuthorizer https://api.monobank.ua/docs/corporate.html
// type CorporateAuthorizer struct {
// 	keyID     string // X-Key-Id
// 	requestID string // X-Request-Id
// 	sign      string // X-Sign (X-time | X-Request-Id | URL)
// }
