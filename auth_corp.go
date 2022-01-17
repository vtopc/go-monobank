package monobank

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
)

// Permissions.
const (
	// PermSt - statements(transactions) and client info of individual(фізичної особи).
	PermSt = "s"
	// PermPI - personal information(first and last names).
	PermPI = "p"
	// PermFOP - statements(transactions) and client info of private entrepreneur(ФОП).
	PermFOP = "f"
)

// Errors.
var (
	ErrDecodePrivateKey  = errors.New("failed to decode private key")
	ErrEncodePublicKey   = errors.New("failed to encode public key with sha1")
	ErrNoPrivateKey      = errors.New("failed to find private key block")
	ErrInvalidEC         = errors.New("invalid elliptic curve private key value")
	ErrInvalidPrivateKey = errors.New("invalid private key length")
)

type CorpAuthMaker struct {
	privateKey *ecdsa.PrivateKey
	KeyID      string // X-Key-Id - ID key of the service
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

const (
	ecPrivateKeyBlockType = "EC PRIVATE KEY"
	ecPrivateKeyVersion   = 1
)

func NewCorpAuthMaker(secKey []byte) (*CorpAuthMaker, error) {
	privateKey, err := decodePrivateKey(secKey)
	if err != nil {
		return nil, ErrDecodePrivateKey
	}

	publicKey := privateKey.PublicKey
	data := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	hash := sha1.New()
	if _, err := hash.Write(data); err != nil {
		return nil, ErrEncodePublicKey
	}
	keyID := hex.EncodeToString(hash.Sum(nil))

	return &CorpAuthMaker{
		privateKey: privateKey,
		KeyID:      keyID,
	}, nil
}

func (c *CorpAuthMaker) New(requestID string) Authorizer {
	return CorpAuth{
		CorpAuthMaker: c,
		requestID:     requestID,
	}
}

func (c *CorpAuthMaker) NewPermissions(permissions ...string) Authorizer {
	return CorpAuth{
		CorpAuthMaker: c,
		permissions:   strings.Join(permissions, ""),
	}
}

type CorpAuth struct {
	*CorpAuthMaker
	requestID   string // Request ID(tokenRequestId)
	permissions string // Permissions
}

func (a CorpAuth) SetAuth(r *http.Request) error {
	if r == nil {
		return nil
	}

	var actor string
	switch {
	case len(a.requestID) > 0:
		actor = a.requestID
		r.Header.Set("X-Request-Id", actor)
	case len(a.permissions) > 0:
		actor = a.permissions
		r.Header.Set("X-Permissions", actor)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	sign, err := a.sign(timestamp, actor, r.URL.Path)
	if err != nil {
		return fmt.Errorf("calculate Sign: %w", err)
	}

	r.Header.Set("X-Key-Id", a.KeyID)
	r.Header.Set("X-Time", timestamp)
	r.Header.Set("X-Sign", sign)

	return nil
}

// sign - calculates Sign (X-time | X-Request-Id/X-Permissions | URL)
func (a CorpAuth) sign(timestamp, actor, urlPath string) (string, error) {
	return a.signString(timestamp + actor + urlPath)
}

func (a CorpAuth) signString(str string) (string, error) {
	hash := sha256.Sum256([]byte(str))

	r, s, err := ecdsa.Sign(rand.Reader, a.privateKey, hash[:])
	if err != nil {
		return "", err
	}

	asn1Data := []*big.Int{r, s}

	bb, err := asn1.Marshal(asn1Data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bb), nil
}

// decodePrivateKey decodes private key into Elliptic Curve Digital Signature Algorithm private key.
// TODO: switch to secp256k1.PrivKeyFromBytes() or https://github.com/square/go-jose/pull/278/files
func decodePrivateKey(b []byte) (*ecdsa.PrivateKey, error) {
	for {
		var privateKeyPemBlock *pem.Block
		privateKeyPemBlock, b = pem.Decode(b)
		if privateKeyPemBlock == nil {
			break
		}

		if privateKeyPemBlock.Type == ecPrivateKeyBlockType {
			ret, err := parseECPrivateKey(privateKeyPemBlock.Bytes)
			if err != nil {
				return nil, err
			}

			return ret, err
		}
	}

	return nil, ErrNoPrivateKey
}

// parseECPrivateKey returns Elliptic Curve Digital Signature Algorithm private key from file content.
// TODO: looks as copy-paste from x509.parseECPrivateKey()
func parseECPrivateKey(b []byte) (*ecdsa.PrivateKey, error) {
	var privKey ecPrivateKey
	if _, err := asn1.Unmarshal(b, &privKey); err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}
	if privKey.Version != ecPrivateKeyVersion {
		//nolint:goerr113
		return nil, fmt.Errorf("unknown EC private key version %d", privKey.Version)
	}

	curve := secp256k1.S256()
	k := new(big.Int).SetBytes(privKey.PrivateKey)
	curveOrder := curve.Params().N
	if k.Cmp(curveOrder) >= 0 {
		return nil, ErrInvalidEC
	}

	priv := new(ecdsa.PrivateKey)
	priv.Curve = curve
	priv.D = k

	privateKey := make([]byte, (curveOrder.BitLen()+7)/8)

	// Some private keys have leading zero padding. This is invalid
	// according to [SEC1], but this code will ignore it.
	for len(privKey.PrivateKey) > len(privateKey) {
		if privKey.PrivateKey[0] != 0 {
			return nil, ErrInvalidPrivateKey
		}
		privKey.PrivateKey = privKey.PrivateKey[1:]
	}

	// Some private keys remove all leading zeros, this is also invalid
	// according to [SEC1] but since OpenSSL used to do this, we ignore
	// this too.
	copy(privateKey[len(privateKey)-len(privKey.PrivateKey):], privKey.PrivateKey)
	priv.X, priv.Y = curve.ScalarBaseMult(privateKey)

	return priv, nil
}
