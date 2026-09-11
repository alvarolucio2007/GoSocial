package auth

import (
	"crypto/rand"
	"errors"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Authenticator interface {
	CreateToken(claims Claims) (string, error)
	VerifyToken(token string) (*Claims, error)
}
type Claims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

func NewClaims(email string, duration time.Duration) (*Claims, error) {
	tokenID := uuid.New()
	now := time.Now()
	payload := &Claims{
		ID:    tokenID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}
	return payload, nil
}

func (c *Claims) Valid() error {
	if time.Now().After(c.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}

func GenerateSymmetricKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}
