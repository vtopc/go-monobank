package monobank

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vtopc/go-rest"
)

// checks that Client satisfies interface
func TestClient_CorporateClient(_ *testing.T) {
	var _ CorporateAPI = CorporateClient{}
}

func TestCorporateClient_withAuth(t *testing.T) {
	authMaker := &CorpAuthMaker{
		KeyID: "abc",
	}

	c, err := NewCorporateClient(nil, authMaker)
	require.NoError(t, err)

	auth := CorpAuth{
		requestID: "123",
	}

	authClient := c.withAuth(auth)
	assert.Equal(t, auth, authClient.auth)
	assert.Equal(t, authMaker, authClient.authMaker)
}

// stubAuth is a no-op Authorizer used by tests that don't care about signing.
type stubAuth struct{}

func (stubAuth) SetAuth(_ *http.Request) error { return nil }

// stubAuthMaker satisfies CorpAuthMakerAPI without doing any real signing.
type stubAuthMaker struct{}

func (stubAuthMaker) New(_ string) Authorizer               { return stubAuth{} }
func (stubAuthMaker) NewPermissions(_ ...string) Authorizer { return stubAuth{} }

func TestCorporateClient_RegistrationStatus(t *testing.T) {
	const pubkeyPEM = "-----BEGIN PUBLIC KEY-----\nABC\n-----END PUBLIC KEY-----\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		hasSuffix(t, r.URL.String(), "/personal/auth/registration/status")

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var got RegistrationStatusRequest
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, []byte(pubkeyPEM), got.Pubkey)

		resp, err := json.Marshal(RegistrationStatusResponse{
			Status: RegistrationStatusApproved,
			KeyID:  "abc123",
		})
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}))
	defer server.Close()

	c, err := NewCorporateClient(nil, stubAuthMaker{})
	require.NoError(t, err)
	c.baseURL = server.URL
	c.restClient = rest.NewClient(server.Client())

	resp, err := c.RegistrationStatus(context.Background(), []byte(pubkeyPEM))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, RegistrationStatusApproved, resp.Status)
	assert.Equal(t, "abc123", resp.KeyID)
}
