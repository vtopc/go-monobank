package monobank

// TODO: add HTTP retry

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const (
	baseURL = "https://api.monobank.ua"
)

type Client struct {
	c       *http.Client
	auth    Authorizer
	baseURL string
}

// New - returns new monobank Client
func New(client *http.Client) Client {
	if client == nil {
		// defaults
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return Client{
		c:       client,
		auth:    NewPublicAuthorizer(),
		baseURL: baseURL,
	}
}

// WithAuth returns copy of Client with authorizer
func (c Client) WithAuth(auth Authorizer) Client {
	c.auth = auth
	return c
}

// WithAuth returns copy of Client with overridden baseURL
func (c Client) WithBaseURL(uri string) Client {
	c.baseURL = uri
	return c
}

// do does request.
// Stores JSON response in the value pointed to by v.
// TODO: make expectedStatusCode a slice:
func (c Client) do(ctx context.Context, req *http.Request, v interface{}, expectedStatusCode int) error {
	// TODO: check that `v` is a pointer or nil

	if req == nil {
		return errors.New("empty request")
	}

	var err error
	req.URL, err = url.Parse(c.baseURL + req.URL.String())
	if err != nil {
		return errors.Wrap(err, "failed to build URL")
	}

	req = req.WithContext(ctx)

	if c.auth != nil {
		c.auth.SetAuth(req)
	}

	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to %s %s", req.Method, req.URL)
	}

	defer resp.Body.Close()

	var body []byte
	if v != nil {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "couldn't read the body")
		}
	}

	switch resp.StatusCode {
	case expectedStatusCode:
		if v == nil {
			// nothing to unmarshal
			return nil
		}

		err = json.Unmarshal(body, v)
		if err == nil {
			return nil
		}

		err = errors.Wrap(err, "failed to unmarshal")

	default:
		err = errors.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return errors.Wrapf(err, "errorBody: %s", string(body))
}
