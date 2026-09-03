package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	Memory      = 64 * 1024
	Iterations  = 3
	Parallelism = 2
	SaltLength  = 16
	KeyLength   = 32
	format      = "$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hashedBytes := argon2.IDKey([]byte(password), salt, Iterations, Memory, Parallelism, KeyLength)

	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hashedBytes)

	fullHash := fmt.Sprintf(format, Memory, Iterations, Parallelism, saltBase64, hashBase64)
	return fullHash, nil
}

func CheckPassword(password, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	originalHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		Iterations,
		Memory,
		Parallelism,
		KeyLength,
	)
	if subtle.ConstantTimeCompare(computedHash, originalHash) == 1 {
		return true, nil
	}
	return false, nil
}
