package store

import "context"

type MockRoleRepository struct {
	roles map[int]*Role
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, ErrNotFound
}
