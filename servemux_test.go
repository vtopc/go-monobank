//go:build go1.22

package monobank

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vtopc/go-rest"
)

// These tests demonstrate Go 1.22 http.ServeMux method+pattern routing and
// r.PathValue() against this library's clients. Compared to the existing
// hasSuffix-based assertions they are strictly stronger:
//
//   - a request to a wrong path returns 404 from the mux (rather than being
//     silently accepted by a catch-all HandlerFunc);
//   - URL segments are extracted by name instead of being recovered from a
//     formatted string.
//
// The file is build-tagged for go1.22 so this PR does not raise the module's
// `go 1.16` minimum; CI runs Go 1.23.x and 1.24.x so the tests execute.

func TestClient_Currency_servemuxPattern(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /bank/currency", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	c := Client{baseURL: server.URL, restClient: rest.NewClient(server.Client())}

	out, err := c.Currency(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestCommonClient_Transactions_pathValues(t *testing.T) {
	var gotAccount, gotFrom, gotTo string

	mux := http.NewServeMux()
	mux.HandleFunc("GET /personal/statement/{accountID}/{from}/{to}",
		func(w http.ResponseWriter, r *http.Request) {
			gotAccount = r.PathValue("accountID")
			gotFrom = r.PathValue("from")
			gotTo = r.PathValue("to")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[]`))
		})
	server := httptest.NewServer(mux)
	defer server.Close()

	cc := newCommonClient(server.Client())
	cc.baseURL = server.URL
	cc.restClient = rest.NewClient(server.Client())

	from := time.Unix(1700000000, 0)
	to := time.Unix(1700001000, 0)
	_, err := cc.Transactions(context.Background(), "acc-123", from, to)
	require.NoError(t, err)

	assert.Equal(t, "acc-123", gotAccount)
	assert.Equal(t, strconv.FormatInt(from.Unix(), 10), gotFrom)
	assert.Equal(t, strconv.FormatInt(to.Unix(), 10), gotTo)
}
