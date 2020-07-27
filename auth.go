package monobank

import (
	"net/http"
)

type Authorizer interface {
	// SetAuth modifies http.Request and sets authorization tokens
	SetAuth(*http.Request) error
}

type PublicAuthorizer struct{}

func (a PublicAuthorizer) SetAuth(_ *http.Request) error {
	return nil
}

func NewPublicAuthorizer() PublicAuthorizer {
	return PublicAuthorizer{}
}
