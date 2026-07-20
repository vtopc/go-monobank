package monobank

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
)

type PublicAPI interface {
	// Currency https://api.monobank.ua/docs/#/definitions/CurrencyInfo
	Currency(context.Context) (Currencies, error)

	// ServerKey fetches the bank's current public key, its identifier and
	// server time. Use it together with VerifyWebhookSignature.
	// https://api.monobank.ua/docs/#tag/Publichni-dani/paths/~1bank~1sync/get
	ServerKey(context.Context) (*ServerKey, error)
}

func (c Client) Currency(ctx context.Context) (Currencies, error) {
	const urlPath = "/bank/currency"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var v Currencies
	err = c.do(req, &v, http.StatusOK)

	return v, err
}

// secp256k1 uncompressed point encoding (SEC 1, §2.3.3):
// 1 prefix byte (0x04) followed by the X and Y coordinates (32 bytes each).
const (
	uncompressedPointPrefix  = 0x04
	secp256k1CoordinateBytes = 32
	uncompressedPointLength  = 1 + 2*secp256k1CoordinateBytes
)

// ErrInvalidPubKey is returned by [Client.ServerKey] when /bank/sync returns
// a serverPubKey that isn't a valid uncompressed secp256k1 point.
var ErrInvalidPubKey = errors.New("invalid serverPubKey: not an uncompressed secp256k1 point")

// bankSyncResponse mirrors the JSON shape of /bank/sync. Kept package-private
// because the public surface is [ServerKey] (built via asServerKey).
type bankSyncResponse struct {
	ServerKeyID    string `json:"serverKeyId"`
	ServerPubKey   string `json:"serverPubKey"`
	ServerTimeMsec int64  `json:"serverTimeMsec"`
}

func (r bankSyncResponse) asServerKey() (*ServerKey, error) {
	pubBytes, err := base64.StdEncoding.DecodeString(r.ServerPubKey)
	if err != nil {
		return nil, fmt.Errorf("decode serverPubKey: %w", err)
	}
	if len(pubBytes) != uncompressedPointLength || pubBytes[0] != uncompressedPointPrefix {
		return nil, ErrInvalidPubKey
	}
	return &ServerKey{
		ID: r.ServerKeyID,
		PubKey: &ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     new(big.Int).SetBytes(pubBytes[1 : 1+secp256k1CoordinateBytes]),
			Y:     new(big.Int).SetBytes(pubBytes[1+secp256k1CoordinateBytes:]),
		},
		ServerTime: time.UnixMilli(r.ServerTimeMsec),
	}, nil
}

// ServerKey fetches the bank's public key used to sign webhooks.
//
// The /bank/sync endpoint is public — no token is required. Cache the result
// and refresh whenever an incoming X-Key-Id stops matching ServerKey.ID; that
// is mono's signal that it rotated the key.
func (c Client) ServerKey(ctx context.Context) (*ServerKey, error) {
	const urlPath = "/bank/sync"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var raw bankSyncResponse
	if err := c.do(req, &raw, http.StatusOK); err != nil {
		return nil, err
	}
	return raw.asServerKey()
}
