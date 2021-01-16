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
	baseURL        = "https://api.monobank.ua"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	httpClient *http.Client
	auth       Authorizer
	baseURL    string // TODO: switch to url.URL
}

// TODO: add WithOpts

// NewClient - returns public monobank Client
func NewClient(client *http.Client) Client {
	if client == nil {
		// defaults
		client = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return Client{
		httpClient: client,
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

	if c.auth != nil { // TODO: return an error if not
		err = c.auth.SetAuth(req)
		if err != nil {
			return NewReqError(req, errors.Wrap(err, "SetAuth"))
		}
	}

	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	err = func() error {
		resp, e := c.httpClient.Do(req)
		if e != nil {
			return e
		}

		defer resp.Body.Close()

		var body []byte
		if v != nil {
			body, e = ioutil.ReadAll(resp.Body)
			if e != nil {
				return errors.Wrap(e, "couldn't read the body")
			}
		}

		// TODO: switch to "for" for multi-status:
		if resp.StatusCode == expectedStatusCode {
			if v == nil {
				// nothing to unmarshal
				return nil
			}

			e = json.Unmarshal(body, v)
			if e == nil {
				return nil
			}

			return errors.Wrap(e, "failed to unmarshal the response body")
		}

		// otherwise, non expected status code:
		return &APIError{
			ResponseStatusCode:  resp.StatusCode,
			ExpectedStatusCodes: []int{expectedStatusCode},
			Err:                 errors.New(string(body)),
		}
	}()
	if err != nil {
		return NewReqError(req, err)
	}

	return nil
}
