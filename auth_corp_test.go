package monobank

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	secKey = []byte("-----BEGIN EC PARAMETERS-----\n" +
		"BgUrgQQACg==\n" +
		"-----END EC PARAMETERS-----\n" +
		"-----BEGIN EC PRIVATE KEY-----\n" +
		"MHQCAQEEIP5DyqGW1yUD5YZRSzsvjT5I9M1utN9aYi3uWJgKhsvPoAcGBSuBBAAK\n" +
		"oUQDQgAEOX+BUepYysBoGR3l9ZsnIXNBm4FYD6m76rGPvbJnUD11xm/SQrOALZYC\n" +
		"s0VrWcLTP60Z1xeLw+NP+D+rUK5IsA==\n" +
		"-----END EC PRIVATE KEY-----\n")

	keyID = "b38daf14d0e6f487949cefbccce99d8add909685"
)

func TestNewCorpAuthMaker(t *testing.T) {
	tests := map[string]struct {
		secKey    []byte
		wantKeyID string
		wantErr   error
	}{
		"positive": {
			secKey:    secKey,
			wantKeyID: keyID,
		},
		"negative": {
			secKey: []byte("invalid"),
			//nolint:goerr113
			wantErr: errors.New("failed to decode private key"),
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			got, err := NewCorpAuthMaker(tt.secKey)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.NotNil(t, got.privateKey)
				assert.Equal(t, tt.wantKeyID, got.KeyID)
			}
		})
	}
}

// ################################## Suite ###################################

type CorpAuthSuite struct {
	suite.Suite
	authMaker *CorpAuthMaker
}

func TestCorpAuthSuite(t *testing.T) {
	suite.Run(t, new(CorpAuthSuite))
}

func (s *CorpAuthSuite) SetupTest() {
	authMaker, err := NewCorpAuthMaker(secKey)
	s.Require().NoError(err)

	s.authMaker = authMaker
}

func (s *CorpAuthSuite) TearDownTest() {
	s.authMaker = nil
}

func (s *CorpAuthSuite) Test_sign() {
	tests := map[string]struct {
		timestamp string
		actor     string
		urlPath   string

		wantErr error
	}{
		"positive": {
			timestamp: "1136239445",
			actor:     "p",
			urlPath:   "/personal/auth/request",
		},
	}

	authIface := s.authMaker.NewPermissions(PermPI)
	auth, ok := authIface.(CorpAuth)
	s.Require().True(ok)

	for name, tt := range tests {
		tt := tt

		s.Run(name, func() {
			got, err := auth.sign(tt.timestamp, tt.actor, tt.urlPath)
			if tt.wantErr != nil {
				s.Assert().EqualError(err, tt.wantErr.Error())
			} else {
				s.Require().NoError(err)
				// result is random, e.g. "MEUCIQCWqCFyU4dYU8RnAhZST7BXW+2giBEu125X/h9lcRnvPwIgU89Dox807jv07XknxRTsgZFlRUsiKUV9KwdiLQ7gWCI="
				s.Assert().Len(got, 96)
			}
		})
	}
}
