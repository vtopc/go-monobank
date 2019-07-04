package monobank

// TODO: add HTTP retry

import (
	"net/http"
	"time"
)

const (
	baseURL = "https://api.monobank.ua"
)

type Client struct {
	c    *http.Client
	auth Authorizer
}

// New - returns new monobank Client
// TODO: replace auth with opts?
func New(client *http.Client, auth Authorizer) Client {
	c := Client{c: client, auth: auth}

	if c.auth == nil {
		c.auth = NewNoopAuthorizer()
	}

	if c.c == nil {
		c.c = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return c
}

// TODO: add Client.NewRequest with user agent
//  req.Header.Set("User-Agent", userAgent)
