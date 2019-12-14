package monobank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type PublicAPI interface {
	// Currency https://api.monobank.ua/docs/#/definitions/CurrencyInfo
	Currency(context.Context) (Currencies, error)
}

type PersonalAPI interface {
	PublicAPI

	// ClientInfo - https://api.monobank.ua/docs/#/definitions/UserInfo
	ClientInfo(context.Context) (*ClientInfo, error)

	// Transactions - gets bank account statements(transations)
	// https://api.monobank.ua/docs/#/definitions/StatementItems
	Transactions(ctx context.Context, accountID string, from, to time.Time) (Transactions, error)

	// SetWebHook - sets webhook for statements
	SetWebHook(ctx context.Context, uri string) error
}

// checks that Client satisfies interface
// TODO: move to test?
var _ PersonalAPI = Client{}

func (c Client) Currency(ctx context.Context) (Currencies, error) {
	const urlSuffix = "/bank/currency"

	req, err := http.NewRequest(http.MethodGet, urlSuffix, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	var v Currencies
	err = c.do(ctx, req, &v, http.StatusOK)
	return v, err
}

func (c Client) ClientInfo(ctx context.Context) (*ClientInfo, error) {
	const urlSuffix = "/personal/client-info"

	req, err := http.NewRequest(http.MethodGet, urlSuffix, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	var v ClientInfo
	err = c.do(ctx, req, &v, http.StatusOK)
	return &v, err
}

// TODO: make `to` optional
func (c Client) Transactions(ctx context.Context, accountID string, from, to time.Time) (
	Transactions, error) {

	const urlSuffix = "/personal/statement"
	uri := fmt.Sprintf("%s/%s/%d/%d", urlSuffix, accountID, from.Unix(), to.Unix())

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	var v Transactions
	err = c.do(ctx, req, &v, http.StatusOK)
	return v, err
}

func (c Client) SetWebHook(ctx context.Context, uri string) error {
	const urlSuffix = "/personal/webhook"

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(WebHookRequest{WebHookURL: uri})
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	req, err := http.NewRequest(http.MethodPost, urlSuffix, &buf)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	return c.do(ctx, req, nil, http.StatusOK)
}
