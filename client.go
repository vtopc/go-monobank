package monobank

// TODO: add HTTP retry

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/vtopc/go-rest"
	"github.com/vtopc/go-rest/defaults"
	"github.com/vtopc/go-rest/interceptors"
)

const (
	baseURL = "https://api.monobank.ua"
)

var ErrEmptyRequest = errors.New("empty request")

type Client struct {
	restClient *rest.Client
	auth       Authorizer
	baseURL    string // TODO: switch to url.URL
}

// TODO: add WithOpts

// NewClient - returns public monobank Client
func NewClient(client *http.Client) Client {
	if client == nil {
		client = defaults.NewHTTPClient()
	}

	_ = interceptors.SetReqContentType(client, "application/json")
	c := rest.NewClient(client)

	return Client{
		restClient: c,
		auth:       NewPublicAuthorizer(),
		baseURL:    baseURL,
	}
}

// WithBaseURL updates baseURL
func (c *Client) WithBaseURL(uri string) {
	c.baseURL = uri
}

func (c *Client) withAuth(auth Authorizer) {
	c.auth = auth
}

// do does request.
// Stores JSON response in the value pointed to by v.
// TODO: make expectedStatusCode a slice:
func (c Client) do(req *http.Request, v interface{}, expectedStatusCode int) error {
	if req == nil {
		return ErrEmptyRequest
	}

	var err error
	req.URL, err = url.Parse(c.baseURL + req.URL.String())
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	if c.auth != nil { // TODO: return an error if not
		err = c.auth.SetAuth(req)
		if err != nil {
			return fmt.Errorf("SetAuth: %w", err)
		}
	}

	err = c.restClient.Do(req, v, expectedStatusCode)
	if err != nil {
		return err
	}

	return nil
}
