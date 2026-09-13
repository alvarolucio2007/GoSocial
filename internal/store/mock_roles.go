package store

import "context"

type MockRoleRepository struct {
	roles map[int]*Role
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	return nil, nil
}
