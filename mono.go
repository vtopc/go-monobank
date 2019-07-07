package monobank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type PublicIface interface {
	// Currency https://api.monobank.ua/docs/#/definitions/CurrencyInfo
	Currency() ([]byte, error)
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

func (c Client) Currency() ([]byte, error) {
	const uri = baseURL + "/bank/currency"

	resp, err := c.c.Get(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to GET: %s", uri)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read response body; status(%d)", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status(%d), body(%s)", resp.StatusCode, string(body))
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		// defaulting to ugly JSON
		return body, nil
	}

	return prettyJSON.Bytes(), nil
}

func (c Client) ClientInfo() (*ClientInfo, error) {
	const urlPrefix = baseURL + "/personal/client-info"

	req, err := http.NewRequest(http.MethodGet, urlPrefix, nil)
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
	const urlPrefix = baseURL + "/personal/statement"
	uri := fmt.Sprintf("%s/%s/%d/%d", urlPrefix, accountID, from.Unix(), to.Unix())

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
