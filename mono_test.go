package monobank

import (
	"testing"
)

// checks that Client satisfies interface
func TestClient_Currency(t *testing.T) {
	var _ PersonalAPI = Client{}
}
