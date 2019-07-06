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
func New(client *http.Client) Client {
	c := Client{c: client, auth: NewNoopAuthorizer()}

	if c.c == nil {
		// defaults
		c.c = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return c
}

// WithAuth returns copy of Client with authorizer
func (c Client) WithAuth(auth Authorizer) Client {
	c.auth = auth
	return c
}

// TODO: add Client.NewRequest with user agent
//  req.Header.Set("User-Agent", userAgent)
