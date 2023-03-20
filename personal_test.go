package monobank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// checks that Client satisfies interface
func TestClient_PersonalAPI(_ *testing.T) {
	var _ PersonalAPI = PersonalClient{}
}

func TestPersonalClient_WithAuth(t *testing.T) {
	c := NewPersonalClient(nil)

	auth := PersAuth{
		token: "123",
	}

	authClient := c.WithAuth(auth)
	assert.Equal(t, auth, authClient.auth)
}
