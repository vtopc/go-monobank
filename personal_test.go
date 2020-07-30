package monobank_test

import (
	"testing"

	"github.com/vtopc/go-monobank"
)

// checks that Client satisfies interface
func TestClient_PersonalAPI(t *testing.T) {
	var _ monobank.PersonalAPI = monobank.PersonalClient{}
}
