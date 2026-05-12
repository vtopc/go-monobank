package monobank

import (
	"context"
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

func TestClient_ServerKey(t *testing.T) {
	const responseBody = `{"serverKeyId":"2626ff34473bb66260b930af946fa9641a06bcd4","serverPubKey":"BNDZP+AGoRC+ER1plDSUCHOw2/aBNIocmD2gS/v34/b0iQ1HBo+oS3/f402e3OXA5uCxakSjuxGMP6X0XP9VIUk=","serverTimeMsec":1778612022098}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		hasSuffix(t, r.URL.String(), "/bank/sync")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	}))
	defer server.Close()

	c := Client{
		baseURL:    server.URL,
		restClient: rest.NewClient(server.Client()),
	}

	sk, err := c.ServerKey(context.Background())
	require.NoError(t, err)
	require.NotNil(t, sk)

	assert.Equal(t, "2626ff34473bb66260b930af946fa9641a06bcd4", sk.ID)
	require.NotNil(t, sk.PubKey)
	require.NotNil(t, sk.PubKey.X)
	require.NotNil(t, sk.PubKey.Y)
	assert.Equal(t, int64(1778612022098), sk.ServerTime.UnixMilli())
}

func TestClient_ServerKey_invalidPubKey(t *testing.T) {
	// serverPubKey decodes to fewer than 65 bytes — bank protocol violation.
	const responseBody = `{"serverKeyId":"x","serverPubKey":"AAAA","serverTimeMsec":0}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	}))
	defer server.Close()

	c := Client{
		baseURL:    server.URL,
		restClient: rest.NewClient(server.Client()),
	}

	_, err := c.ServerKey(context.Background())
	assert.ErrorIs(t, err, ErrInvalidPubKey)
}
