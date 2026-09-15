package store

import (
	"context"
)

type FollowKey struct {
	UserID     int64
	FollowerID int64
}
type MockFollowerRepository struct {
	followers map[FollowKey]struct{}
}

func (m *MockFollowerRepository) Follow(ctx context.Context, userID, followerID int64) error {
	followKey := FollowKey{UserID: userID, FollowerID: followerID}
	if _, exist := m.followers[followKey]; exist {
		return ErrConflict
	}
	m.followers[followKey] = struct{}{}
	return nil
}

func (m *MockFollowerRepository) Unfollow(ctx context.Context, userID, followerID int64) error {
	followKey := FollowKey{UserID: userID, FollowerID: followerID}
	if _, exist := m.followers[followKey]; !exist {
		return ErrNotFound
	}
	delete(m.followers, followKey)
	return nil
}
