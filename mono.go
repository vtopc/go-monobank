package monobank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

	// SetWebHook - sets webhook for statements
	SetWebHook(uri string) error
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

	var v Statements
	err = c.Do(req, http.StatusOK, &v)
	return v, err
}

func (c Client) SetWebHook(uri string) error {
	const urlSuffix = "/personal/webhook"

	var buf *bytes.Buffer
	err := json.NewEncoder(buf).Encode(webHook{WebHookURL: uri})
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	req, err := http.NewRequest(http.MethodPost, urlSuffix, buf)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	return c.Do(req, http.StatusOK, nil)
}
