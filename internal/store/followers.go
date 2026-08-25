package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}
type FollowerStore struct {
	db *sql.DB
}
type FollowerRepository interface {
	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `DELETE followers WHERE user_id=$1 AND follower_id=$2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}
