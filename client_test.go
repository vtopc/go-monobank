package monobank

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_do(t *testing.T) {
	tests := map[string]struct {
		method             string
		urlPostfix         string
		expectedStatusCode int
		body               []byte
		v                  Transaction
		want               Transaction
	}{
		"positive-get": {
			method:             http.MethodGet,
			urlPostfix:         "/test",
			expectedStatusCode: http.StatusOK,
			body:               []byte(`{"description":"test"}`),
			v:                  Transaction{},
			want:               Transaction{Description: "test"},
		},
	}

	for k, tc := range tests {
		tc := tc
		t.Run(k, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.method, r.Method)
				hasSuffix(t, r.URL.String(), tc.urlPostfix)

				// Send response to be tested
				w.WriteHeader(tc.expectedStatusCode)
				_, _ = w.Write(tc.body)
			}))
			defer server.Close()

			c := Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			req, err := http.NewRequest(tc.method, tc.urlPostfix, http.NoBody)
			require.NoError(t, err)

			// test:
			err = c.do(context.Background(), req, &tc.v, tc.expectedStatusCode)
			require.NoError(t, err)
			assert.Equal(t, tc.want, tc.v)
		})
	}
}

func hasSuffix(t *testing.T, s, suffix string) {
	assert.Truef(t, strings.HasSuffix(s, suffix), "expected '%s' to ends with suffix '%s'", s, suffix)
}

func TestClient_withAuth(t *testing.T) {
	c := NewClient(nil)

	auth := PersAuth{
		token: "123",
	}

	c.withAuth(auth)
	assert.Equal(t, auth, c.auth)
}
