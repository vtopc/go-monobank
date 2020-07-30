package monobank

import (
	"context"
	"net/http"
	"time"
)

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

type PersonalClient struct {
	commonClient
}

func NewPersonalClient(client *http.Client) PersonalClient {
	return PersonalClient{
		commonClient: newCommonClient(client),
	}
}

func (c PersonalClient) SetWebHook(ctx context.Context, uri string) error {
	const urlPath = "/personal/webhook"

	return c.setWebHook(ctx, uri, urlPath)
}

// WithAuth returns copy of PersonalClient with authorizer
func (c PersonalClient) WithAuth(auth Authorizer) PersonalClient {
	c.withAuth(auth)
	return c
}
