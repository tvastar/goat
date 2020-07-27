package goat

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/pkg/errors"
)

type pkceInfo struct {
	verifier, challenge, method string
}

func newPKCEInfo() (*pkceInfo, error) {
	// Per https://tools.ietf.org/html/rfc7636#section-4.1
	data := make([]byte, 64)
	n, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, errors.New("insufficient randomness")
	}

	verifier := base64.RawURLEncoding.EncodeToString(data)
	sum := sha256.Sum256([]byte(verifier))
	info := &pkceInfo{
		verifier:  verifier,
		challenge: base64.RawURLEncoding.EncodeToString(sum[:]),
		method:    "S256",
	}
	return info, nil
}
