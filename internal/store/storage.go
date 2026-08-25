package store

import (
	"database/sql"
	"time"
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

const QueryTimeout time.Duration = 5 * time.Second
