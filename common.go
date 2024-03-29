package monobank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CommonAPI interface {
	PublicAPI

	// SetWebHook - sets webhook for statements
	SetWebHook(ctx context.Context, uri string) error
}

// commonClient contains common to Personal and Corporate API
type commonClient struct {
	Client
}

func newCommonClient(client *http.Client) commonClient {
	return commonClient{
		Client: NewClient(client),
	}
}

func (c commonClient) ClientInfo(ctx context.Context) (*ClientInfo, error) {
	const urlPath = "/personal/client-info"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var v ClientInfo
	err = c.do(req, &v, http.StatusOK)

	return &v, err
}

// TODO: make `to` optional
func (c commonClient) Transactions(ctx context.Context, accountID string, from, to time.Time) (
	Transactions, error) {

	const urlPath = "/personal/statement"
	uri := fmt.Sprintf("%s/%s/%d/%d", urlPath, accountID, from.Unix(), to.Unix())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var v Transactions
	err = c.do(req, &v, http.StatusOK)

	return v, err
}

func (c commonClient) setWebHook(ctx context.Context, uri, urlPath string) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(WebHookRequest{WebHookURL: uri})
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.do(req, nil, http.StatusOK)
}
