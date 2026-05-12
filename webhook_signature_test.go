package monobank

import (
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"math/big"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test vector captured from a real mono personal webhook (a jar top-up).
// The serverPubKey was fetched from https://api.monobank.ua/bank/sync at
// the time the webhook arrived; key ID 2626ff34... matches the X-Key-Id
// header that came alongside the body and signature below.
const (
	testServerPubKeyB64 = "BNDZP+AGoRC+ER1plDSUCHOw2/aBNIocmD2gS/v34/b0iQ1HBo+oS3/f402e3OXA5uCxakSjuxGMP6X0XP9VIUk="
	testServerKeyID     = "2626ff34473bb66260b930af946fa9641a06bcd4"

	testWebhookBody = `{"type":"StatementItem","data":{"account":"GdIXP9tJybhRwW4yl457iw","statementItem":{"id":"XfHmJ0KH0p1jVAw58w","time":1778612175,"description":"Поповнення «Донатики🥰»","mcc":4829,"originalMcc":4829,"amount":-20000,"operationAmount":-20000,"currencyCode":980,"commissionRate":0,"cashbackAmount":0,"balance":62360,"hold":true,"receiptId":"ACAB-X504-7TP3-143T"}}}`
	testWebhookSign = "hqiXbiDWs4lCkg/i9cZaqWBHgvio2PhGNnzkBiU7MBJxnODkMf8RYKsaLke8gbm+1XSvZgmMHPYFw3XswwL7Qw=="
)

func testServerPubKey(t *testing.T) *ecdsa.PublicKey {
	t.Helper()
	raw, err := base64.StdEncoding.DecodeString(testServerPubKeyB64)
	require.NoError(t, err)
	require.Equal(t, 65, len(raw))
	require.Equal(t, byte(0x04), raw[0])
	return &ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     new(big.Int).SetBytes(raw[1:33]),
		Y:     new(big.Int).SetBytes(raw[33:]),
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	pub := testServerPubKey(t)

	t.Run("valid signature", func(t *testing.T) {
		err := VerifyWebhookSignature(pub, []byte(testWebhookBody), testWebhookSign)
		assert.NoError(t, err)
	})

	t.Run("tampered body", func(t *testing.T) {
		tampered := []byte(testWebhookBody[:len(testWebhookBody)-2] + "9}")
		err := VerifyWebhookSignature(pub, tampered, testWebhookSign)
		assert.ErrorIs(t, err, ErrBadSignature)
	})

	t.Run("garbage signature", func(t *testing.T) {
		err := VerifyWebhookSignature(pub, []byte(testWebhookBody), "AAAA")
		assert.ErrorIs(t, err, ErrBadSignature)
	})

	t.Run("non-base64 signature", func(t *testing.T) {
		err := VerifyWebhookSignature(pub, []byte(testWebhookBody), "not base64!!!")
		assert.ErrorIs(t, err, ErrBadSignatureEncoding)
	})

	t.Run("nil pubkey", func(t *testing.T) {
		err := VerifyWebhookSignature(nil, []byte(testWebhookBody), testWebhookSign)
		assert.ErrorIs(t, err, ErrMissingPubKey)
	})
}

func TestVerifyWebhookSignature_errorsAreDistinguishable(t *testing.T) {
	// guard against accidentally collapsing sentinel errors
	assert.False(t, errors.Is(ErrBadSignature, ErrBadSignatureEncoding))
	assert.False(t, errors.Is(ErrBadSignatureEncoding, ErrMissingPubKey))
}
