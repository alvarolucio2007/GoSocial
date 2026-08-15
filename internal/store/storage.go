package store

import (
	"database/sql"
)

type Storage struct {
	PostRepository
	UserRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		PostRepository: &PostStore{db: db},
		UserRepository: &UserStore{db: db},
	}
}
