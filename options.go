package monobank

import (
	"net/http"

	"github.com/vtopc/go-rest"
	"github.com/vtopc/go-rest/interceptors"
)

// Option configures a [Client]. Pass options to [New] to override defaults.
// Designed to be additive — future options (HTTP retry, timeouts, custom
// interceptors) slot in without breaking existing callers.
type Option func(*Client)

// WithHTTPClient sets a custom *http.Client. If not used, the client returned
// by [defaults.NewHTTPClient] from github.com/vtopc/go-rest is used.
//
// Note: monobank requires the Content-Type interceptor; if you bring your own
// client it will be configured automatically.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient == nil {
			return
		}
		_ = interceptors.SetReqContentType(httpClient, "application/json")
		c.restClient = rest.NewClient(httpClient)
	}
}

// WithBaseURL overrides the default https://api.monobank.ua base URL.
// Useful for testing against a recorded server or against monobank's sandbox.
func WithBaseURL(uri string) Option {
	return func(c *Client) {
		c.baseURL = uri
	}
}

// New returns a public [Client] built from the supplied options. This is the
// extensible counterpart to [NewClient] — both produce equivalent clients
// when given equivalent inputs.
//
//	c := monobank.New(
//	    monobank.WithHTTPClient(myHTTP),
//	    monobank.WithBaseURL("https://staging.example/"),
//	)
func New(opts ...Option) Client {
	c := NewClient(nil)
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
