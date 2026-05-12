package monobank

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
)

// Webhook signature verification errors.
var (
	// ErrBadSignatureEncoding is returned when X-Sign cannot be base64-decoded.
	ErrBadSignatureEncoding = errors.New("X-Sign is not valid base64")
	// ErrMissingPubKey is returned when a nil public key is passed in.
	ErrMissingPubKey = errors.New("missing public key")
)

// VerifyWebhookSignature returns nil iff the X-Sign header is a valid ECDSA
// signature of body produced by the bank's serverPubKey.
//
// Usage:
//
//	sk, err := client.ServerKey(ctx)
//	...
//	if err := monobank.VerifyWebhookSignature(sk.PubKey, body, r.Header.Get("X-Sign")); err != nil {
//	    // reject
//	}
//
// Cache sk and refresh it when the incoming X-Key-Id header stops matching
// sk.ID — that is mono's signal that the key rotated.
//
// The bank currently sends the signature as a raw 64-byte r||s pair encoded
// in base64; ASN.1 DER (RFC 5480) is accepted as well, in case mono ever
// switches encodings.
func VerifyWebhookSignature(pub *ecdsa.PublicKey, body []byte, xSign string) error {
	if pub == nil {
		return ErrMissingPubKey
	}
	sig, err := base64.StdEncoding.DecodeString(xSign)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBadSignatureEncoding, err)
	}

	digest := sha256.Sum256(body)

	if len(sig) == 64 {
		r := new(big.Int).SetBytes(sig[:32])
		s := new(big.Int).SetBytes(sig[32:])
		if ecdsa.Verify(pub, digest[:], r, s) {
			return nil
		}
	}

	var asn1Sig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(sig, &asn1Sig); err == nil &&
		asn1Sig.R != nil && asn1Sig.S != nil &&
		ecdsa.Verify(pub, digest[:], asn1Sig.R, asn1Sig.S) {
		return nil
	}

	return ErrBadSignature
}

// ErrBadSignature is returned by VerifyWebhookSignature when the signature
// does not match. Compare with errors.Is.
var ErrBadSignature = errors.New("webhook signature is invalid")
