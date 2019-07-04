package monobank

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Iface interface {
	// Currency https://api.monobank.ua/docs/#/definitions/CurrencyInfo
	Currency() ([]byte, error)

	// https://api.monobank.ua/docs/#/definitions/UserInfo
	ClientInfo() (*ClientInfo, error)
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
	const uri = baseURL + "/personal/client-info"

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	c.auth.SetAuth(req)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to GET: %s", uri)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var v ClientInfo
		err = json.NewDecoder(resp.Body).Decode(&v)
		if err == nil {
			return &v, nil
		}

		err = errors.Wrap(err, "failed to unmarshal")

	default:
		err = errors.Errorf("unexpected status(%d)", resp.StatusCode)
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		return nil, errors.Wrapf(err, "failed to read response body: %s", bodyErr)
	}

	return nil, errors.Wrapf(err, "body: %s", string(body))
}
