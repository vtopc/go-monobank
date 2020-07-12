package monobank

import (
	"testing"
)

// checks that Client satisfies interface
func TestClient_PersonalAPI(t *testing.T) {
	var _ PersonalAPI = Client{}
}
