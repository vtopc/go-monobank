package monobank

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReqError(t *testing.T) {
	uri, err := url.Parse("http://example.com/call")
	require.NoError(t, err)

	tests := map[string]struct {
		err        error
		wantError  string
		wantUnwrap string
	}{
		"short": {
			err: &ReqError{
				Method: http.MethodGet,
				URL:    uri,
				Err:    errors.New("SetAuth"),
			},
			wantError:  "request GET http://example.com/call: SetAuth",
			wantUnwrap: "SetAuth",
		},
		"nested": {
			err: &ReqError{
				Method: http.MethodGet,
				URL:    uri,
				Err: &APIError{
					ResponseStatusCode:  400,
					ExpectedStatusCodes: []int{200},
					Err:                 errors.New(`{"errorDescription":"42"}`),
				},
			},
			wantError:  "request GET http://example.com/call: unexpected status code 400(want [200]): unmarshal",
			wantUnwrap: `unexpected status code 400(want [200]): {"errorDescription":"42"}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run("Error()", func(t *testing.T) {
				assert.EqualError(t, tt.err, tt.wantError)
			})

			t.Run("Unwrap()", func(t *testing.T) {
				got := errors.Unwrap(tt.err)
				assert.EqualError(t, got, tt.wantUnwrap)
			})
		})
	}
}
