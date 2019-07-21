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

type PublicAuthorizer struct{}

func (a PublicAuthorizer) SetAuth(req *http.Request) {}

func NewPublicAuthorizer() PublicAuthorizer {
	return PublicAuthorizer{}
}

// TODO: add NewCorporateAuthorizer https://api.monobank.ua/docs/corporate.html
// type CorporateAuthorizer struct {
// 	keyID     string // X-Key-Id
// 	requestID string // X-Request-Id
// 	sign      string // X-Sign (X-time | X-Request-Id | URL)
// }
//
// Sign (X-time | X-Request-Id | URL)
//  - X-time=1561560962 (Wed, 26 Jun 2019 14:56:02 GMT)
//  - X-Request-Id=acW5k2ERnupgnWFyBkCY0nA
//  - URL=/personal/client-info
//  Sign (1561560962acW5k2ERnupgnWFyBkCY0nA/personal/client-info)=eaBHn+T18kr7w6uAzLJ037o1w/JMpAfV81yNaHaXJUxv9bbi/cORuA0gGsazwG+VxCq2Y+TmIb81zbGbuiaRQA==
