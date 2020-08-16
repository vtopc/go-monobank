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
		v         *ReqError
		wantError string
		wantCause string
	}{
		"unmarshal": {
			v: &ReqError{
				Method: http.MethodGet,
				URL:    uri,
				Err:    errors.New("SetAuth"),
			},
			wantError: "request GET http://example.com/call: SetAuth",
			wantCause: "SetAuth",
		},
		"nested": {
			v: &ReqError{
				Method: http.MethodGet,
				URL:    uri,
				Err: &APIError{
					ResponseStatusCode:  400,
					ExpectedStatusCodes: []int{200},
					Err:                 errors.New("unmarshal"),
				},
			},
			wantError: "request GET http://example.com/call: unexpected status code 400(want [200]): unmarshal",
			wantCause: "unexpected status code 400(want [200]): unmarshal",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run("Error()", func(t *testing.T) {
				assert.Equal(t, tt.v.Error(), tt.wantError)
			})

			t.Run("Cause()", func(t *testing.T) {
				assert.EqualError(t, tt.v.Cause(), tt.wantCause)
			})
		})
	}
}
