package monobank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// checks that Client satisfies interface
func TestClient_PersonalAPI(t *testing.T) {
	var _ PersonalAPI = PersonalClient{}
}

func TestPersonalClient_WithAuth(t *testing.T) {
	c := PersonalClient{
		commonClient: newCommonClient(nil),
	}

	auth := PersAuth{
		token: "123",
	}

	c.WithAuth(auth)
	assert.Equal(t, auth, c.auth)
}
