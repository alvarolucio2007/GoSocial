package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const QueryTimeout time.Duration = 5 * time.Second

var (
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type Storage struct {
	Posts     PostRepository
	Users     UserRepository
	Comments  CommentRepository
	Followers FollowerRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		if err2 := tx.Rollback(); err != nil {
			return err2
		}
		return err
	}
	return tx.Commit()
}
