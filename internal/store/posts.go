package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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
	User      User      `json:"user"`
}
type PostWithMetadata struct {
	Post
	CommentCount int `json:"comments_count"`
}
type PostStore struct {
	db *sql.DB
}
type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	Read(ctx context.Context, idPost int) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, idPost int) error
	GetUserFeed(ctx context.Context, idUser int64) ([]PostWithMetadata, error)
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content,title,user_id,tags)
									VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	tagsStr := "{}"
	if len(post.Tags) > 0 {
		tagsStr = fmt.Sprintf("{%s}", strings.Join(post.Tags, ","))
	}
	if err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, tagsStr).
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

func (s *PostStore) GetUserFeed(ctx context.Context, idUser int64) ([]PostWithMetadata, error) {
	query := ` 
	SELECT 
		p.id,p.user_id,p.title,p.content,p.created_at,p.version,p.tags, u.username,
		COUNT(c.id) AS comments_count
	FROM posts p
	LEFT JOIN comments c on c.post_id = p.id
	left join users u on p.user_id = u.id
	join followers f on f.follower_id = p.user_id  or p.user_id = $1
	where f.user_id = $1 or p.user_id = $1
	GROUP BY p.id,u.username
	ORDER BY p.created_at DESC;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		var tagsBytes []byte
		err := rows.Scan(
			&p.ID, &p.UserID,
			&p.Title, &p.Content,
			&p.CreatedAt, &p.Version,
			&tagsBytes, &p.User.Username, &p.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		p.Tags = parseTags(tagsBytes)
		feed = append(feed, p)
	}
	return feed, nil
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
