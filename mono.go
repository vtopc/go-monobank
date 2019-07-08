package monobank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type PublicIface interface {
	// Currency https://api.monobank.ua/docs/#/definitions/CurrencyInfo
	Currency() (Currencies, error)
}

type Iface interface {
	PublicIface

	// ClientInfo - https://api.monobank.ua/docs/#/definitions/UserInfo
	ClientInfo() (*ClientInfo, error)

	// Statement - bank account statement(transations)
	// https://api.monobank.ua/docs/#/definitions/StatementItems
	Statement(accountID string, from, to time.Time) (Statements, error)
}

// checks that Client satisfies interface
// TODO: move to test?
var _ Iface = Client{}

func (c Client) Currency() (Currencies, error) {
	const urlSuffix = "/bank/currency"

	req, err := http.NewRequest(http.MethodGet, urlSuffix, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	var v Currencies
	err = c.Do(req, http.StatusOK, &v)
	return v, err
}

func (c Client) ClientInfo() (*ClientInfo, error) {
	const urlSuffix = "/personal/client-info"

	req, err := http.NewRequest(http.MethodGet, urlSuffix, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	c.auth.SetAuth(req)

	var v ClientInfo
	err = c.Do(req, http.StatusOK, &v)
	return &v, err
}

// TODO: make `to` optional
func (c Client) Statement(accountID string, from, to time.Time) (Statements, error) {
	const urlSuffix = "/personal/statement"
	uri := fmt.Sprintf("%s/%s/%d/%d", urlSuffix, accountID, from.Unix(), to.Unix())

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	c.auth.SetAuth(req)

	var v Statements
	err = c.Do(req, http.StatusOK, &v)
	return v, err
}

func (c Client) Do(req *http.Request, statusCode int, v interface{}) error {
	// TODO: check that `v` is a pointer

	var err error
	req.URL, err = url.Parse(baseURL + req.URL.String())
	if err != nil {
		return errors.Wrap(err, "failed to parse URL")
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
