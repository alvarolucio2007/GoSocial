package store

import "context"

type MockFollowerRepository struct {
	followers map[int]*Follower
}

func (m *MockFollowerRepository) Follow(context.Context, int64, int64) error {
	return nil
}

func (m *MockFollowerRepository) Unfollow(context.Context, int64, int64) error {
	return nil
}
