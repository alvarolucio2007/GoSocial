package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
}
type PostStore struct {
	db *sql.DB
}
type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	Read(ctx context.Context, idPost int) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, idPost int) error
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content,title,user_id,tags)
									VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	tags := pgtype.Array[string]{
		Elements: []string{},
		Valid:    true,
	}
	if post.Tags != nil {
		tags.Elements = post.Tags
	} else {
		post.Tags = []string{}
	}
	if err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, tags).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return err
	}
	return nil
}

var ErrPostNotFound = errors.New("post not found")

func (s *PostStore) Read(ctx context.Context, idPost int) (*Post, error) {
	query := `SELECT id,content,title,user_id,tags,created_at,updated_at ,version FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	var p Post
	var tagsBytes []byte
	err := s.db.QueryRowContext(ctx, query, idPost).Scan(&p.ID, &p.Content, &p.Title, &p.UserID, &tagsBytes, &p.CreatedAt, &p.UpdatedAt, &p.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}
	p.Tags = parseTags(tagsBytes)
	return &p, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts
	SET
		title = COALESCE(NULLIF($1, ''), title, ''),
		content = COALESCE(NULLIF($2, ''), content, ''),
		tags = COALESCE(NULLIF($3::text[], '{}'::text[]), tags),
		updated_at = $4, version=version+1
	WHERE id = $5 AND version = $6
	RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	updateTime := time.Now()
	log.Println(post.Tags)
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.Tags, updateTime, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrPostNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostStore) Delete(ctx context.Context, idPost int) error {
	query := `DELETE FROM posts WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, idPost)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrPostNotFound
	}
	return nil
}

func parseTags(b []byte) []string {
	if len(b) == 0 || string(b) == "{}" || string(b) == "NULL" {
		return []string{}
	}
	s := strings.Trim(string(b), "{}")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
