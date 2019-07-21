package monobank

// TODO: add HTTP retry

import (
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
	c    *http.Client
	auth Authorizer
}

// New - returns new monobank Client
func New(client *http.Client) Client {
	c := Client{c: client, auth: NewPublicAuthorizer()}

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

// Do does request.
// `statusCode` - expected HTTP status code from response.
func (c Client) Do(req *http.Request, statusCode int, v interface{}) error {
	// TODO: check that `v` is a pointer or nil

	var err error
	req.URL, err = url.Parse(baseURL + req.URL.String())
	if err != nil {
		return errors.Wrap(err, "failed to build URL")
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to %s %s", req.Method, req.URL)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case statusCode:
		if v == nil {
			// nothing to unmarshal
			return nil
		}

		err = json.NewDecoder(resp.Body).Decode(&v)
		if err == nil {
			return nil
		}

		err = errors.Wrap(err, "failed to unmarshal")

	default:
		err = errors.Errorf("unexpected status(%d)", resp.StatusCode)
	}

	errorBody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return errors.Wrapf(err, "but failed to read response body: %s", e)
	}

	return errors.Wrapf(err, "errorBody: %s", string(errorBody))
}
