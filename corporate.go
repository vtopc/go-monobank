package monobank

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrEmptyAuthMaker = errors.New("authMaker is nil")

type CorporateAPI interface {
	CommonAPI

	// Settings
	// https://api.monobank.ua/docs/corporate.html#tag/Avtorizaciya-ta-nalashtuvannya-kompaniyi/paths/~1personal~1corp~1settings/get
	Settings(ctx context.Context) (*CorpSettings, error)

	// Auth initializes client access.
	// https://api.monobank.ua/docs/corporate.html#tag/Kliyentski-personalni-dani/paths/~1personal~1auth~1request/post
	Auth(ctx context.Context, callbackURL string, permissions ...string) (*TokenRequest, error)

	// CheckAuth checks status of request for client's personal data.
	// https://api.monobank.ua/docs/corporate.html#tag/Kliyentski-personalni-dani/paths/~1personal~1auth~1request/get
	CheckAuth(ctx context.Context, requestID string) error

	// ClientInfo
	// https://api.monobank.ua/docs/corporate.html#tag/Kliyentski-personalni-dani/paths/~1personal~1client-info/get
	ClientInfo(ctx context.Context, requestID string) (*ClientInfo, error)

	// Transactions - gets bank account statements(transactions)
	// https://api.monobank.ua/docs/corporate.html#tag/Kliyentski-personalni-dani/paths/~1personal~1statement~1{account}~1{from}~1{to}/get
	Transactions(ctx context.Context, requestID, accountID string, from, to time.Time) (Transactions, error)
}

type CorpAuthMakerAPI interface {
	// New returns corp Authorizer for endpoints with Request ID.
	New(requestID string) Authorizer

	// NewPermissions returns corp Authorizer for Auth endpoint to get Request ID.
	// Omitting permissions means all permissions.
	NewPermissions(permissions ...string) Authorizer
}

type CorporateClient struct {
	commonClient
	authMaker CorpAuthMakerAPI
}

const urlPathAuth = "/personal/auth/request"

// NewCorporateClient returns corporate client
func NewCorporateClient(client *http.Client, authMaker CorpAuthMakerAPI) (CorporateClient, error) {
	if authMaker == nil {
		return CorporateClient{}, ErrEmptyAuthMaker
	}

	return CorporateClient{
		commonClient: newCommonClient(client),
		authMaker:    authMaker,
	}, nil
}

// Auth initializes access.
func (c CorporateClient) Auth(ctx context.Context, callbackURL string, permissions ...string) (*TokenRequest, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPathAuth, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Callback", callbackURL)

	authClient := c.withAuth(c.authMaker.NewPermissions(permissions...))

	var v TokenRequest
	err = authClient.commonClient.do(req, &v, http.StatusOK)

	return &v, err
}

func (c CorporateClient) CheckAuth(ctx context.Context, requestID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPathAuth, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	authClient := c.withAuth(c.authMaker.New(requestID))

	return authClient.do(req, nil, http.StatusOK)
}

// SetWebHook sets webhook for corporate API.
func (c CorporateClient) SetWebHook(ctx context.Context, uri string) error {
	const urlPath = "/personal/corp/webhook"

	authClient := c.withAuth(c.authMaker.New(""))

	return authClient.setWebHook(ctx, uri, urlPath)
}

func (c CorporateClient) ClientInfo(ctx context.Context, requestID string) (*ClientInfo, error) {
	authClient := c.withAuth(c.authMaker.New(requestID))

	return authClient.commonClient.ClientInfo(ctx)
}

func (c CorporateClient) Transactions(ctx context.Context, requestID, accountID string, from, to time.Time) (Transactions, error) {
	authClient := c.withAuth(c.authMaker.New(requestID))

	return authClient.commonClient.Transactions(ctx, accountID, from, to)
}

// Settings
// https://api.monobank.ua/docs/corporate.html#tag/Avtorizaciya-ta-nalashtuvannya-kompaniyi/paths/~1personal~1corp~1settings/get
func (c CorporateClient) Settings(ctx context.Context) (*CorpSettings, error) {
	const urlPath = "/personal/corp/settings"

	authClient := c.withAuth(c.authMaker.New(""))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var v CorpSettings
	err = authClient.do(req, &v, http.StatusOK)

	return &v, err
}

// withAuth returns copy of CorporateClient with authorizer
// TODO: remove?
func (c CorporateClient) withAuth(auth Authorizer) CorporateClient {
	c.commonClient.withAuth(auth)

	return c
}
