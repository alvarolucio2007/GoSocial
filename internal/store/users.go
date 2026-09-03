package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alvarolucio2007/GoSocial/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}
type UserStore struct {
	db *sql.DB
}
type UserRepository interface {
	Create(context.Context, *User) error
	Read(context.Context, int) (*User, error)
	Update(context.Context, *User) error
	Delete(context.Context, int) error
}
type password struct {
	text string
	hash string
}

func (p *password) Set(text string) error {
	hash, err := utils.HashPassword(text)
	if err != nil {
		return err
	}
	p.text = text
	p.hash = hash
	return nil
}

func (p *password) Compare(text string) (bool, error) {
	isSame, err := utils.CheckPassword(text, p.hash)
	if err != nil {
		return false, err
	}
	return isSame, nil
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username,email,password)
									VALUES ($1,$2,$3) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	if err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).Scan(&user.ID, &user.CreatedAt); err != nil {
		return err
	}
	return nil
}

var ErrUserNotFound = errors.New("user not found")

func (s *UserStore) Read(ctx context.Context, idUser int) (*User, error) {
	query := `SELECT id,username,email,password,created_at FROM users WHERE id = $1`
	var u User

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, idUser).Scan(&u.ID, &u.Username, &u.Email, &u.Password.hash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
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
		password=COALESCE(NULLIF($3,'\x'::bytea),password)
	WHERE id=$4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, user.Username, user.Email, user.Password.hash, user.ID)
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
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
