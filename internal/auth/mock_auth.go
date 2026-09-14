package auth

import "uuid"

type MockAuthenticator struct{ tokens map[string]*Claims }

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{
		tokens: make(map[string]*Claims),
	}
}

func (m *MockAuthenticator) CreateToken(claims Claims) (string, error) {
	token := uuid.New().String()
	m.tokens[token] = &claims
	return token, nil
}

func (m *MockAuthenticator) VerifyToken(token string) (*Claims, error) {
	claims, ok := m.tokens[token]
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
