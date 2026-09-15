package store

import (
	"context"
	"errors"
)

type FollowKey struct {
	UserID     int64
	FollowerID int64
}
type MockFollowerRepository struct {
	followers map[FollowKey]struct{}
}

func (m *MockFollowerRepository) Follow(ctx context.Context, userID int64, followerID int64) error {
	followKey := FollowKey{UserID: userID, FollowerID: followerID}
	m.followers[followKey] = struct{}{}
	return nil
}

func (m *MockFollowerRepository) Unfollow(ctx context.Context, userID int64, followerID int64) error {
	followKey := FollowKey{UserID: userID, FollowerID: followerID}
	if _, exist := m.followers[followKey]; !exist {
		return errors.New("placeholder")
	}
	delete(m.followers[followKey])
}
