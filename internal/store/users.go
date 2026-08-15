package store

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}
type UserStore struct {
	db *sql.DB
}
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Read(ctx context.Context, idUser int) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, idUser int) error
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username,email,password)
									VALUES ($1,$2,$3) RETURNING id,created_at`
	if err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt); err != nil {
		return err
	}
	return nil
}

var ErrUserNotFound = errors.New("user not found")

func (s *UserStore) Read(ctx context.Context, idUser int) (*User, error) {
	query := `SELECT id,username,email,password,created_at FROM usuario WHERE id = $1`
	var u User
	err := s.db.QueryRowContext(ctx, query, idUser).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
	if errors.Is(err, ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `UPDATE users
	SET	
		username=COALESCE(NULLIF($1,''),username),
		email=COALESCE(NULLIF($2,''),email),
		password=COALESCE(NULLIF($3,''),password)
	WHERE id=$4
	`
	res, err := s.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *UserStore) Delete(ctx context.Context, idUser int) error {
	query := `DELETE FROM users WHERE id=$1`
	res, err := s.db.ExecContext(ctx, query, idUser)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrUserNotFound
	}
	return nil
}
