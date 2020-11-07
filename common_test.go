package monobank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_commonClient_withAuth(t *testing.T) {
	c := newCommonClient(nil)

	auth := PersAuth{
		token: "123",
	}

	c.withAuth(auth)
	assert.Equal(t, auth, c.auth)
}
