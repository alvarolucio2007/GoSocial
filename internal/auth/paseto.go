package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/o1egl/paseto"
)

type PasetoAuthenticator struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

var ErrInvalidKeySize = errors.New("invalid symmetric key size: must be exactly 32 characters")

func NewPasetoAuthenticator(symmetricKey []byte) (*PasetoAuthenticator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	auth := &PasetoAuthenticator{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}
	return auth, nil
}

func (p *PasetoAuthenticator) CreateToken(claims Claims) (string, error) {
	return p.paseto.Encrypt(p.symmetricKey, claims, nil)
}

func (p *PasetoAuthenticator) VerifyToken(token string) (*Claims, error) {
	claims := &Claims{}
	err := p.paseto.Decrypt(token, p.symmetricKey, claims, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = claims.Valid()
	if err != nil {
		return nil, err
	}
	return claims, nil
}
