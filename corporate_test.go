package monobank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// checks that Client satisfies interface
func TestClient_CorporateClient(t *testing.T) {
	var _ CorporateAPI = CorporateClient{}
}

func TestCorporateClient_withAuth(t *testing.T) {
	c := CorporateClient{
		commonClient: newCommonClient(nil),
	}

	auth := CorpAuth{
		requestID: "123",
	}

	authClient := c.withAuth(auth)
	assert.Equal(t, auth, authClient.auth)
}
