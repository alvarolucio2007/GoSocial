package auth

type MockAuthenticator struct{}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{}
}

func (m *MockAuthenticator) CreateToken(claims Claims) (string, error) {
	return "", nil
}

func (m *MockAuthenticator) VerifyToken(token string) (*Claims, error) {
	return nil, nil
}
