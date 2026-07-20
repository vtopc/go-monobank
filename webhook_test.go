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

// Test vector captured from a real mono personal-API webhook (a jar top-up).
// The serverPubKey was fetched from /bank/sync at the time the webhook
// arrived; serverKeyId matches the X-Key-Id header that accompanied the
// payload and signature below.
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
		assert.NoError(t, VerifyWebhookSignature(pub, []byte(testWebhookBody), testWebhookSign))
	})

	t.Run("tampered body", func(t *testing.T) {
		tampered := []byte(testWebhookBody[:len(testWebhookBody)-2] + "9}")
		assert.ErrorIs(t, VerifyWebhookSignature(pub, tampered, testWebhookSign), ErrBadSignature)
	})

	t.Run("garbage signature", func(t *testing.T) {
		assert.ErrorIs(t, VerifyWebhookSignature(pub, []byte(testWebhookBody), "AAAA"), ErrBadSignature)
	})

	t.Run("non-base64 signature", func(t *testing.T) {
		assert.ErrorIs(t, VerifyWebhookSignature(pub, []byte(testWebhookBody), "not base64!!!"),
			ErrBadSignatureEncoding)
	})

	t.Run("nil pubkey", func(t *testing.T) {
		assert.ErrorIs(t, VerifyWebhookSignature(nil, []byte(testWebhookBody), testWebhookSign),
			ErrMissingPubKey)
	})
}

func TestServerKey_Verify(t *testing.T) {
	sk := &ServerKey{ID: testServerKeyID, PubKey: testServerPubKey(t)}

	assert.NoError(t, sk.Verify([]byte(testWebhookBody), testWebhookSign))
	assert.ErrorIs(t, sk.Verify([]byte("tampered"), testWebhookSign), ErrBadSignature)

	var nilSK *ServerKey
	assert.ErrorIs(t, nilSK.Verify([]byte(testWebhookBody), testWebhookSign), ErrMissingPubKey)
}

func TestParseWebHook(t *testing.T) {
	t.Run("known type", func(t *testing.T) {
		w, err := ParseWebHook([]byte(testWebhookBody))
		require.NoError(t, err)
		require.NotNil(t, w)
		assert.Equal(t, WebHookTypeStatementItem, w.Type)
		assert.Equal(t, "GdIXP9tJybhRwW4yl457iw", w.Data.AccountID)
		assert.Equal(t, "XfHmJ0KH0p1jVAw58w", w.Data.Transaction.ID)
		assert.Equal(t, int64(-20000), w.Data.Transaction.Amount)
		assert.Equal(t, "ACAB-X504-7TP3-143T", w.Data.Transaction.ReceiptID)
	})

	t.Run("unknown type still parsed, error flagged", func(t *testing.T) {
		w, err := ParseWebHook([]byte(`{"type":"NewEventKind","data":{}}`))
		assert.ErrorIs(t, err, ErrUnknownWebHookType)
		require.NotNil(t, w)
		assert.Equal(t, "NewEventKind", w.Type)
	})

	t.Run("malformed JSON", func(t *testing.T) {
		_, err := ParseWebHook([]byte(`{`))
		require.Error(t, err)
		assert.False(t, errors.Is(err, ErrUnknownWebHookType))
	})
}
