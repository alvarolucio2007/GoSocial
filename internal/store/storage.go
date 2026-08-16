package store

import (
	"database/sql"
)

type Storage struct {
	Posts PostRepository
	Users UserRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db: db},
		Users: &UserStore{db: db},
	}
}
