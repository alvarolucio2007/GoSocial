package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}
type UserStore struct {
	db *sql.DB
}
type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *User) error
	Read(ctx context.Context, userID int) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID int) error
	CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
	Activate(ctx context.Context, token string) error
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

func mapPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		switch pgErr.ConstraintName {
		case "users_email_key":
			return ErrDuplicateEmail
		case "users_username_key":
			return ErrDuplicateUsername
		}
		return ErrConflict
	default:
		return err
	}
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (username,email,password)
									VALUES ($1,$2,$3) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	if err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).Scan(&user.ID, &user.CreatedAt); err != nil {
		return mapPgError(err)
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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token,user_id,expiry) VALUES ($1,$2,$3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `SELECT u.id,u.username,u.email,u.created_at,u.is_active
	FROM users u JOIN user_invitations ui
	ON u.id=ui.user_id
	WHERE ui.token =  $1 AND ui.expiry > $2`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	user := &User{}
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	if err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrPostNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id=$1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}
