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

// TODO: replace auth with opts?
func New(client *http.Client, auth Authorizer) Client {
	c := Client{auth: auth}

	if client != nil {
		c.c = client
	} else {
		c.c = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return c
}

// TODO: add Client.NewRequest with user agent
//  req.Header.Set("User-Agent", userAgent)
