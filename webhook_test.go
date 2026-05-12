package monobank

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestServerKey_Verify(t *testing.T) {
	pub := testServerPubKey(t)
	sk := &ServerKey{ID: testServerKeyID, PubKey: pub}

	assert.NoError(t, sk.Verify([]byte(testWebhookBody), testWebhookSign))
	assert.ErrorIs(t, sk.Verify([]byte("tampered"), testWebhookSign), ErrBadSignature)

	var nilSK *ServerKey
	assert.ErrorIs(t, nilSK.Verify([]byte(testWebhookBody), testWebhookSign), ErrMissingPubKey)
}
