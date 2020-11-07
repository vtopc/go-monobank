package monobank

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// checks that Client satisfies interface
func TestClient_CorporateClient(t *testing.T) {
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
