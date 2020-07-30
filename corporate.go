package monobank

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type CorporateAPI interface {
	PersonalAPI

	// Auth initializes access.
	//  Use corporate authorizer for it
	Auth(ctx context.Context, callbackURL string) (*TokenRequest, error)

	// CheckAuth checks status of request for client's personal data.
	CheckAuth(context.Context) error
}

type CorporateClient struct {
	commonClient
}

const urlPathAuth = "/personal/auth/request"

func NewCorporateClient(client *http.Client) CorporateClient {
	return CorporateClient{
		commonClient: newCommonClient(client),
	}
}

func (c CorporateClient) Auth(ctx context.Context, callbackURL string) (*TokenRequest, error) {
	req, err := http.NewRequest(http.MethodPost, urlPathAuth, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("X-Callback", callbackURL)

	var v TokenRequest
	err = c.do(ctx, req, &v, http.StatusOK)

	return &v, err
}

// CheckAuth checks status of request for client's personal data.
func (c CorporateClient) CheckAuth(ctx context.Context) error {
	req, err := http.NewRequest(http.MethodGet, urlPathAuth, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	return c.do(ctx, req, nil, http.StatusOK)
}

func (c CorporateClient) SetWebHook(ctx context.Context, uri string) error {
	const urlPath = "/personal/corp/webhook"

	return c.setWebHook(ctx, uri, urlPath)
}

// WithAuth returns copy of CorporateClient with authorizer
func (c CorporateClient) WithAuth(auth Authorizer) CorporateClient {
	c.withAuth(auth)
	return c
}
