package monobank

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Known webhook payload types.
const (
	// WebHookTypeStatementItem is the only "type" mono currently sends for
	// personal-API webhooks (a single bank-account statement entry).
	WebHookTypeStatementItem = "StatementItem"
)

// ErrUnknownWebHookType is returned by ParseWebHook when the payload's "type"
// field is neither empty nor a known constant. It is not necessarily an error
// for callers — mono may add new event types — but it lets you opt out of
// processing payloads you don't recognise.
var ErrUnknownWebHookType = errors.New("unknown webhook type")

// ParseWebHook decodes a raw webhook body into a *WebHookResponse and rejects
// payloads whose top-level "type" is unknown. Verify the signature with
// VerifyWebhookSignature before calling this — ParseWebHook does no auth.
func ParseWebHook(body []byte) (*WebHookResponse, error) {
	var v WebHookResponse
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, fmt.Errorf("decode webhook: %w", err)
	}
	switch v.Type {
	case WebHookTypeStatementItem:
		return &v, nil
	default:
		return &v, fmt.Errorf("%w: %q", ErrUnknownWebHookType, v.Type)
	}
}
