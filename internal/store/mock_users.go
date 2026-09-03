package store

import (
	"context"
	"database/sql"
	"time"
)

type MockUserRepository struct {
	users map[int]*User
}

func (m *MockUserRepository) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	m.users[int(user.ID)] = user
	return nil
}

func (m *MockUserRepository) Read(ctx context.Context, id int) (*User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *User) error {
	oldUser, err := m.Read(ctx, int(user.ID))
	if err != nil {
		return ErrUserNotFound
	}
	u := *oldUser
	if user.Username != "" {
		u.Username = user.Username
	}
	if user.Email != "" {
		u.Email = user.Email
	}
	password := password{}
	if err := password.Set(""); err != nil {
		return err
	}
	if user.Password != password {
		u.Password = user.Password
	}
	m.users[int(user.ID)] = &u
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, idUser int) error {
	_, err := m.Read(ctx, idUser)
	if err != nil {
		return ErrPostNotFound
	}
	delete(m.users, idUser)
	return nil
}

func (m *MockUserRepository) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}
