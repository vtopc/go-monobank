package monobank

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vtopc/go-rest"
)

// checks that Client satisfies PublicAPI
func TestClient_PublicAPI(_ *testing.T) {
	var _ PublicAPI = Client{}
}

// validServerKey is a captured /bank/sync response used by ServerKey tests.
// The pubkey is the production one at the time of capture; we only assert
// shape / parsing, so it's safe to keep static.
var validServerKey = bankSyncResponse{
	ServerKeyID:    "2626ff34473bb66260b930af946fa9641a06bcd4",
	ServerPubKey:   "BNDZP+AGoRC+ER1plDSUCHOw2/aBNIocmD2gS/v34/b0iQ1HBo+oS3/f402e3OXA5uCxakSjuxGMP6X0XP9VIUk=",
	ServerTimeMsec: 1778612022098,
}

func TestClient_ServerKey(t *testing.T) {
	body, err := json.Marshal(validServerKey)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		hasSuffix(t, r.URL.String(), "/bank/sync")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	c := Client{
		baseURL:    server.URL,
		restClient: rest.NewClient(server.Client()),
	}

	sk, err := c.ServerKey(context.Background())
	require.NoError(t, err)
	require.NotNil(t, sk)

	assert.Equal(t, validServerKey.ServerKeyID, sk.ID)
	require.NotNil(t, sk.PubKey)
	require.NotNil(t, sk.PubKey.X)
	require.NotNil(t, sk.PubKey.Y)
	assert.Equal(t, validServerKey.ServerTimeMsec, sk.ServerTime.UnixMilli())
}

func TestClient_ServerKey_invalidPubKey(t *testing.T) {
	body, err := json.Marshal(bankSyncResponse{
		ServerKeyID:  "x",
		ServerPubKey: "AAAA", // decodes to <65 bytes — not a valid uncompressed point
	})
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	c := Client{
		baseURL:    server.URL,
		restClient: rest.NewClient(server.Client()),
	}

	_, err = c.ServerKey(context.Background())
	assert.ErrorIs(t, err, ErrInvalidPubKey)
}
